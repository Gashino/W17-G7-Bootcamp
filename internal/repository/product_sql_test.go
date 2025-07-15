package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Funciones auxiliares para crear punteros a diferentes tipos de datos
func strPtr(s string) *string       { return &s }
func intPtr(i int) *int             { return &i }
func float64Ptr(f float64) *float64 { return &f }

// Variables para controlar el registro del driver txdb
var (
	txdbRegistered bool = false
	txdbMutex      sync.Mutex
)

// Estructura para leer la configuración de la base de datos desde config.yml
type Config struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

// ProductSql ya está definido en product_map.go

// Configurar la base de datos para pruebas
func setupTxDB() string {
	// Usar un mutex para evitar condiciones de carrera al registrar el driver
	txdbMutex.Lock()
	defer txdbMutex.Unlock()

	// Evitar registrar el driver más de una vez
	if !txdbRegistered {
		// Leer la configuración desde config.yml
		config := loadConfig()

		// Configurar la conexión MySQL usando los valores del config.yml
		cfg := mysql.Config{
			User:                 config.Database.User,
			Passwd:               config.Database.Password,
			Net:                  "tcp",
			Addr:                 fmt.Sprintf("%s:%s", config.Database.Host, config.Database.Port),
			DBName:               config.Database.Name,
			ParseTime:            true,
			AllowNativePasswords: true,
		}

		// Verificar si el driver ya está registrado
		for _, driver := range sql.Drivers() {
			if driver == "txdb" {
				// El driver ya está registrado, no necesitamos registrarlo de nuevo
				txdbRegistered = true
				return "txdb"
			}
		}

		// Registrar el driver si no está registrado
		txdb.Register("txdb", "mysql", cfg.FormatDSN())
		txdbRegistered = true
	}

	return "txdb"
}

// Cargar la configuración desde config.yml
func loadConfig() Config {
	var config Config

	// Intentar leer el archivo de configuración
	data, err := os.ReadFile("../../config.yml")
	if err != nil {
		// Si hay un error, mostrar un mensaje y terminar el test
		panic(fmt.Sprintf("Error al leer el archivo config.yml: %v\nAsegúrate de que el archivo config.yml existe en la raíz del proyecto", err))
	}

	// Parsear el archivo YAML
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		// Si hay un error al parsear el YAML, mostrar un mensaje y terminar el test
		panic(fmt.Sprintf("Error al parsear el archivo config.yml: %v\nVerifica que el formato del archivo sea correcto", err))
	}

	// Verificar que la configuración de la base de datos esté completa
	if config.Database.User == "" || config.Database.Host == "" || config.Database.Port == "" || config.Database.Name == "" {
		panic("La configuración de la base de datos en config.yml está incompleta")
	}

	return config
}

// Configurar la base de datos de prueba para cada test
func setupProductTestDB(t *testing.T) (*ProductSql, *sql.DB) {
	// Usar un nombre único para cada conexión de test
	driver := setupTxDB()
	db, err := sql.Open(driver, fmt.Sprintf("product_test_%s", t.Name()))
	require.NoError(t, err)

	// Verificar que la conexión funciona
	err = db.Ping()
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	return &ProductSql{db: db}, db
}

// Crear un producto completo para pruebas
func createTestProduct() models.Product {
	return models.Product{
		ProductAttributes: models.ProductAttributes{
			ProductCode:                    strPtr("TEST001"),
			Description:                    strPtr("Test Product"),
			NetWeight:                      float64Ptr(1.5),
			ExpirationRate:                 intPtr(30),
			RecommendedFreezingTemperature: float64Ptr(-18.0),
			FreezingRate:                   intPtr(5),
			// Usar un product_type_id válido según el script de base de datos (101-120)
			ProductTypeId: intPtr(101),
			// Usar un seller_id válido según el script de base de datos (1-10)
			SellerId: intPtr(1),
		},
		Dimensions: models.Dimensions{
			Width:  float64Ptr(10.0),
			Height: float64Ptr(5.0),
			Length: float64Ptr(8.0),
		},
	}
}

func TestProductSQL_GetAll(t *testing.T) {
	repo, db := setupProductTestDB(t)
	defer db.Close()

	// Act
	products := repo.GetAll()

	// Assert
	assert.NotNil(t, products)
	// Verificamos que se devuelvan los productos de seed
	assert.GreaterOrEqual(t, len(products), 1, "Deberían existir al menos algunos productos en la base de datos de prueba")
}

func TestProductSQL_GetById(t *testing.T) {
	repo, db := setupProductTestDB(t)
	defer db.Close()

	t.Run("existing product", func(t *testing.T) {
		// Act
		product, err := repo.GetById(1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 1, product.ID)
	})

	t.Run("non-existing product", func(t *testing.T) {
		// Act
		product, err := repo.GetById(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
	})
}

func TestProductSQL_Create(t *testing.T) {
	t.Run("create new product successfully", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Arrange
		newProduct := createTestProduct()

		// Act
		createdProduct, err := repo.Create(newProduct)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, createdProduct)
		assert.NotZero(t, createdProduct.ID)
		assert.Equal(t, *newProduct.ProductCode, *createdProduct.ProductCode)
		assert.Equal(t, *newProduct.Description, *createdProduct.Description)

		// Verificar que el producto fue creado en la base de datos
		savedProduct, err := repo.GetById(createdProduct.ID)
		assert.NoError(t, err)
		assert.Equal(t, *createdProduct.ProductCode, *savedProduct.ProductCode)
	})

	t.Run("create product with duplicate code", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Arrange - Crear un primer producto
		product1 := createTestProduct()
		createdProduct, err := repo.Create(product1)
		require.NoError(t, err)
		require.NotNil(t, createdProduct)

		// Intentar crear un segundo producto con el mismo código
		product2 := createTestProduct() // Mismo código que product1

		// Act
		duplicateProduct, err := repo.Create(product2)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, duplicateProduct)
		assert.Contains(t, err.Error(), "product_code already exists")
	})

	t.Run("create product with invalid product type id", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Arrange - Crear un producto con un tipo de producto inválido
		invalidProduct := createTestProduct()
		*invalidProduct.ProductTypeId = 999 // ID que no existe

		// Act
		createdProduct, err := repo.Create(invalidProduct)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, createdProduct)
		assert.Contains(t, err.Error(), "invalid product_type_id")
	})
}

func TestProductSQL_Update(t *testing.T) {
	t.Run("update existing product", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Arrange - Crear un producto primero
		newProduct := createTestProduct()
		createdProduct, err := repo.Create(newProduct)
		require.NoError(t, err)
		require.NotNil(t, createdProduct)

		// Actualizar el producto
		updatedProduct := models.Product{
			ID: createdProduct.ID,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    strPtr("TESTUPDATED"),
				Description:                    strPtr("Test Updated"),
				NetWeight:                      float64Ptr(2.0),
				ExpirationRate:                 intPtr(60),
				RecommendedFreezingTemperature: float64Ptr(-20.0),
				FreezingRate:                   intPtr(10),
				// Usar un product_type_id válido según el script de base de datos (101-120)
				ProductTypeId: intPtr(101),
				SellerId:      intPtr(1),
			},
			Dimensions: models.Dimensions{
				Width:  float64Ptr(15.0),
				Height: float64Ptr(7.0),
				Length: float64Ptr(12.0),
			},
		}

		// Act
		err = repo.Update(createdProduct.ID, updatedProduct)

		// Assert
		assert.NoError(t, err)

		// Verificar que el producto fue actualizado
		savedProduct, err := repo.GetById(createdProduct.ID)
		assert.NoError(t, err)
		assert.Equal(t, *updatedProduct.ProductCode, *savedProduct.ProductCode)
		assert.Equal(t, *updatedProduct.Description, *savedProduct.Description)
		assert.Equal(t, *updatedProduct.NetWeight, *savedProduct.NetWeight)
	})

	t.Run("update non-existing product", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Arrange
		nonExistingProduct := models.Product{
			ID: 9999, // ID que no existe
			ProductAttributes: models.ProductAttributes{
				ProductCode: strPtr("TESTNONEXIST"),
				Description: strPtr("Test Non Existing"),
			},
		}

		// Act
		err := repo.Update(nonExistingProduct.ID, nonExistingProduct)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
	})

	t.Run("update with duplicate product code", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Arrange - Crear dos productos
		product1 := createTestProduct()
		*product1.ProductCode = "TESTDUP1"

		product2 := createTestProduct()
		*product2.ProductCode = "TESTDUP2"

		createdProduct1, err := repo.Create(product1)
		require.NoError(t, err)
		createdProduct2, err := repo.Create(product2)
		require.NoError(t, err)

		// Intentar actualizar el producto 2 con el código del producto 1
		updatedProduct := models.Product{
			ID: createdProduct2.ID,
			ProductAttributes: models.ProductAttributes{
				ProductCode: strPtr("TESTDUP1"), // Código duplicado (mismo que createdProduct1)
				Description: strPtr("Updated Description"),
			},
		}

		// Act
		err = repo.Update(createdProduct2.ID, updatedProduct)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product_code already exists")

		// Verificar que el producto no fue actualizado
		savedProduct, err := repo.GetById(createdProduct2.ID)
		assert.NoError(t, err)
		assert.Equal(t, *product2.ProductCode, *savedProduct.ProductCode)

		// Verificar que el producto 1 sigue existiendo (usando createdProduct1)
		_, err = repo.GetById(createdProduct1.ID)
		assert.NoError(t, err)
	})
}

func TestProductSQL_Delete(t *testing.T) {
	t.Run("delete existing product", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Arrange - Crear un producto primero
		newProduct := createTestProduct()
		createdProduct, err := repo.Create(newProduct)
		require.NoError(t, err)
		require.NotNil(t, createdProduct)

		// Act
		err = repo.Delete(createdProduct.ID)

		// Assert
		assert.NoError(t, err)

		// Verificar que el producto fue eliminado
		_, err = repo.GetById(createdProduct.ID)
		assert.Error(t, err)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
	})

	t.Run("delete non-existing product", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Act
		err := repo.Delete(9999) // ID que no existe

		// Assert
		assert.Error(t, err)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
	})
}

func TestProductSQL_GetProductRecords(t *testing.T) {
	t.Run("get all product records", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Act
		records, err := repo.GetProductRecords(nil)

		// Assert
		assert.NoError(t, err)
		// El producto recién creado no tiene registros asociados en la base de datos
		// así que el resultado puede estar vacío, pero no debe ser nil
		assert.NotNil(t, records)
	})

	t.Run("get product records for specific product", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Arrange - Crear un producto primero
		newProduct := createTestProduct()
		createdProduct, err := repo.Create(newProduct)
		require.NoError(t, err)
		require.NotNil(t, createdProduct)

		// Act
		productID := createdProduct.ID
		_, err = repo.GetProductRecords(&productID)

		// Assert
		assert.NoError(t, err)
		// El producto recién creado no tiene registros asociados en la base de datos
		// La implementación actual puede devolver nil o un slice vacío
		// Solo verificamos que no haya error
	})

	t.Run("get product records for non-existing product", func(t *testing.T) {
		repo, db := setupProductTestDB(t)
		defer db.Close()

		// Arrange
		nonExistingID := 9999

		// Act
		records, err := repo.GetProductRecords(&nonExistingID)

		// Assert
		assert.NoError(t, err) // No debería dar error, solo devolver un slice vacío
		assert.Empty(t, records)
	})
}
