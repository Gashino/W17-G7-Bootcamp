package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Variables para controlar el registro del driver txdb para seller tests
var (
	sellerTxdbRegistered bool = false
	sellerTxdbMutex      sync.Mutex
)

// Estructura para leer la configuración de la base de datos desde config.yml
type SellerConfig struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

// Configurar la base de datos para pruebas de seller
func setupSellerTxDB() string {
	// Usar un mutex para evitar condiciones de carrera al registrar el driver
	sellerTxdbMutex.Lock()
	defer sellerTxdbMutex.Unlock()

	// Evitar registrar el driver más de una vez
	if !sellerTxdbRegistered {
		// Leer la configuración desde config.yml
		config := loadSellerConfig()

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

		// Registrar el driver txdb
		txdb.Register("txdb_seller", "mysql", cfg.FormatDSN())
		sellerTxdbRegistered = true
	}

	return "txdb_seller"
}

// Cargar la configuración desde config.yml
func loadSellerConfig() SellerConfig {
	var config SellerConfig

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
		panic(fmt.Sprintf("Error al parsear el archivo config.yml: %v", err))
	}

	return config
}

// Función helper para crear una base de datos de prueba
func createTestSellerDB() *sql.DB {
	driverName := setupSellerTxDB()
	db, err := sql.Open(driverName, fmt.Sprintf("seller_test_%d", time.Now().UnixNano()))
	if err != nil {
		panic(fmt.Sprintf("Error al conectar a la base de datos: %v", err))
	}
	return db
}

// Test para Create - Caso exitoso
func TestSellerSql_Create_Success(t *testing.T) {
	db := createTestSellerDB()
	defer db.Close()

	repo := NewSellerSql(db)

	// Crear un seller de prueba
	seller := models.Seller{
		SellerAttributes: models.SellerAttributes{
			CId:         "S001",
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "123-456-7890",
			LocalityID:  1,
		},
	}

	// Ejecutar el test
	result, err := repo.Create(seller)

	// Verificar resultado
	assert.NoError(t, err)
	assert.NotZero(t, result.ID)
	assert.Equal(t, seller.CId, result.CId)
	assert.Equal(t, seller.CompanyName, result.CompanyName)
	assert.Equal(t, seller.Address, result.Address)
	assert.Equal(t, seller.Telephone, result.Telephone)
	assert.Equal(t, seller.LocalityID, result.LocalityID)
}

// Test para Create - Caso de CId duplicado
func TestSellerSql_Create_DuplicateCId(t *testing.T) {
	db := createTestSellerDB()
	defer db.Close()

	repo := NewSellerSql(db)

	// Crear el primer seller
	seller1 := models.Seller{
		SellerAttributes: models.SellerAttributes{
			CId:         "S001",
			CompanyName: "Test Company 1",
			Address:     "123 Test St",
			Telephone:   "123-456-7890",
			LocalityID:  1,
		},
	}

	_, err := repo.Create(seller1)
	require.NoError(t, err)

	// Intentar crear otro seller con el mismo CId
	seller2 := models.Seller{
		SellerAttributes: models.SellerAttributes{
			CId:         "S001", // CId duplicado
			CompanyName: "Test Company 2",
			Address:     "456 Test Ave",
			Telephone:   "987-654-3210",
			LocalityID:  1,
		},
	}

	// Ejecutar el test
	_, err = repo.Create(seller2)

	// Verificar que se devuelve un error de conflicto
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrConflict], err)
}

// Test para GetById - Caso exitoso
func TestSellerSql_GetById_Success(t *testing.T) {
	db := createTestSellerDB()
	defer db.Close()

	repo := NewSellerSql(db)

	// Crear un seller de prueba
	seller := models.Seller{
		SellerAttributes: models.SellerAttributes{
			CId:         "S001",
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "123-456-7890",
			LocalityID:  1,
		},
	}

	created, err := repo.Create(seller)
	require.NoError(t, err)

	// Ejecutar el test
	result, err := repo.GetById(created.ID)

	// Verificar resultado
	assert.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, created.CId, result.CId)
	assert.Equal(t, created.CompanyName, result.CompanyName)
	assert.Equal(t, created.Address, result.Address)
	assert.Equal(t, created.Telephone, result.Telephone)
	assert.Equal(t, created.LocalityID, result.LocalityID)
}

// Test para GetById - Caso de seller no encontrado
func TestSellerSql_GetById_NotFound(t *testing.T) {
	db := createTestSellerDB()
	defer db.Close()

	repo := NewSellerSql(db)

	// Ejecutar el test con un ID que no existe
	_, err := repo.GetById(999)

	// Verificar que se devuelve un error de no encontrado
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
}

// Test para FindAll - Caso exitoso
func TestSellerSql_FindAll_Success(t *testing.T) {
	db := createTestSellerDB()
	defer db.Close()

	repo := NewSellerSql(db)

	// Crear varios sellers de prueba
	sellers := []models.Seller{
		{
			SellerAttributes: models.SellerAttributes{
				CId:         "S001",
				CompanyName: "Test Company 1",
				Address:     "123 Test St",
				Telephone:   "123-456-7890",
				LocalityID:  1,
			},
		},
		{
			SellerAttributes: models.SellerAttributes{
				CId:         "S002",
				CompanyName: "Test Company 2",
				Address:     "456 Test Ave",
				Telephone:   "987-654-3210",
				LocalityID:  2,
			},
		},
	}

	// Crear los sellers
	for _, seller := range sellers {
		_, err := repo.Create(seller)
		require.NoError(t, err)
	}

	// Ejecutar el test
	result, err := repo.FindAll()

	// Verificar resultado
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 2) // At least 2 records (may have seeded data)

	// Check that our test records are in the result
	foundSeller1 := false
	foundSeller2 := false
	for _, seller := range result {
		if seller.CId == sellers[0].CId {
			foundSeller1 = true
		}
		if seller.CId == sellers[1].CId {
			foundSeller2 = true
		}
	}
	assert.True(t, foundSeller1, "Seller 1 should be found in results")
	assert.True(t, foundSeller2, "Seller 2 should be found in results")
}

// Test para UpdateFields - Caso exitoso
func TestSellerSql_UpdateFields_Success(t *testing.T) {
	db := createTestSellerDB()
	defer db.Close()

	repo := NewSellerSql(db)

	// Crear un seller de prueba
	seller := models.Seller{
		SellerAttributes: models.SellerAttributes{
			CId:         "S001",
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "123-456-7890",
			LocalityID:  1,
		},
	}

	created, err := repo.Create(seller)
	require.NoError(t, err)

	// Preparar datos para actualización
	updateData := models.SellerCreateRequest{
		CompanyName: stringPtr("Updated Company"),
		Address:     stringPtr("Updated Address"),
		Telephone:   stringPtr("555-123-4567"),
	}

	// Ejecutar el test
	result, err := repo.UpdateFields(created.ID, updateData)

	// Verificar resultado
	assert.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, created.CId, result.CId) // CId no debería cambiar
	assert.Equal(t, *updateData.CompanyName, result.CompanyName)
	assert.Equal(t, *updateData.Address, result.Address)
	assert.Equal(t, *updateData.Telephone, result.Telephone)
	assert.Equal(t, created.LocalityID, result.LocalityID) // LocalityID no debería cambiar
}

// Test para UpdateFields - Caso de seller no encontrado
func TestSellerSql_UpdateFields_NotFound(t *testing.T) {
	db := createTestSellerDB()
	defer db.Close()

	repo := NewSellerSql(db)

	updateData := models.SellerCreateRequest{
		CompanyName: stringPtr("Updated Company"),
	}

	// Ejecutar el test con un ID que no existe
	_, err := repo.UpdateFields(999, updateData)

	// Verificar que se devuelve un error de no encontrado
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
}

// Test para DeleteSeller - Caso exitoso
func TestSellerSql_DeleteSeller_Success(t *testing.T) {
	db := createTestSellerDB()
	defer db.Close()

	repo := NewSellerSql(db)

	// Crear un seller de prueba
	seller := models.Seller{
		SellerAttributes: models.SellerAttributes{
			CId:         "S001",
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "123-456-7890",
			LocalityID:  1,
		},
	}

	created, err := repo.Create(seller)
	require.NoError(t, err)

	// Ejecutar el test
	err = repo.DeleteSeller(created.ID)

	// Verificar resultado
	assert.NoError(t, err)

	// Verificar que el seller ya no existe
	_, err = repo.GetById(created.ID)
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
}

// Test para DeleteSeller - Caso de seller no encontrado
func TestSellerSql_DeleteSeller_NotFound(t *testing.T) {
	db := createTestSellerDB()
	defer db.Close()

	repo := NewSellerSql(db)

	// Ejecutar el test con un ID que no existe
	err := repo.DeleteSeller(999)

	// Verificar que se devuelve un error de no encontrado
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
}

// Helper function para crear punteros a string
func stringPtr(s string) *string {
	return &s
}
