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

// Variables para controlar el registro del driver txdb para locality tests
var (
	localityTxdbRegistered bool = false
	localityTxdbMutex      sync.Mutex
)

// Estructura para leer la configuración de la base de datos desde config.yml
type LocalityConfig struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

// Configurar la base de datos para pruebas de locality
func setupLocalityTxDB() string {
	// Usar un mutex para evitar condiciones de carrera al registrar el driver
	localityTxdbMutex.Lock()
	defer localityTxdbMutex.Unlock()

	// Evitar registrar el driver más de una vez
	if !localityTxdbRegistered {
		// Leer la configuración desde config.yml
		config := loadLocalityConfig()

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
		txdb.Register("txdb_locality", "mysql", cfg.FormatDSN())
		localityTxdbRegistered = true
	}

	return "txdb_locality"
}

// Cargar la configuración desde config.yml
func loadLocalityConfig() LocalityConfig {
	var config LocalityConfig

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
func createTestLocalityDB() *sql.DB {
	driverName := setupLocalityTxDB()
	db, err := sql.Open(driverName, fmt.Sprintf("locality_test_%d", time.Now().UnixNano()))
	if err != nil {
		panic(fmt.Sprintf("Error al conectar a la base de datos: %v", err))
	}
	return db
}

// Test para Create - Caso exitoso
func TestLocalitySql_Create_Success(t *testing.T) {
	db := createTestLocalityDB()
	defer db.Close()

	repo := NewLocalitySql(db)

	// Crear una locality de prueba
	locality := models.Locality{
		LocalitiesAttributes: models.LocalitiesAttributes{
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		},
	}

	// Ejecutar el test
	result, err := repo.Create(locality)

	// Verificar resultado
	assert.NoError(t, err)
	assert.NotZero(t, result.ID)
	assert.Equal(t, locality.LocalityName, result.LocalityName)
	assert.Equal(t, locality.ProvinceName, result.ProvinceName)
	assert.Equal(t, locality.CountryName, result.CountryName)
}

// Test para Create - Caso de nombre duplicado
func TestLocalitySql_Create_DuplicateName(t *testing.T) {
	db := createTestLocalityDB()
	defer db.Close()

	repo := NewLocalitySql(db)

	// Crear la primera locality
	locality1 := models.Locality{
		LocalitiesAttributes: models.LocalitiesAttributes{
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		},
	}

	_, err := repo.Create(locality1)
	require.NoError(t, err)

	// Intentar crear otra locality con el mismo nombre
	locality2 := models.Locality{
		LocalitiesAttributes: models.LocalitiesAttributes{
			LocalityName: "Buenos Aires", // Nombre duplicado
			ProvinceName: "CABA",
			CountryName:  "Argentina",
		},
	}

	// Ejecutar el test
	_, err = repo.Create(locality2)

	// Verificar que se devuelve un error de conflicto
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrConflict], err)
}

// Test para GetById - Caso exitoso
func TestLocalitySql_GetById_Success(t *testing.T) {
	db := createTestLocalityDB()
	defer db.Close()

	repo := NewLocalitySql(db)

	// Crear una locality de prueba
	locality := models.Locality{
		LocalitiesAttributes: models.LocalitiesAttributes{
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		},
	}

	created, err := repo.Create(locality)
	require.NoError(t, err)

	// Ejecutar el test
	result, err := repo.GetById(created.ID)

	// Verificar resultado
	assert.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, created.LocalityName, result.LocalityName)
	assert.Equal(t, created.ProvinceName, result.ProvinceName)
	assert.Equal(t, created.CountryName, result.CountryName)
}

// Test para GetById - Caso de locality no encontrada
func TestLocalitySql_GetById_NotFound(t *testing.T) {
	db := createTestLocalityDB()
	defer db.Close()

	repo := NewLocalitySql(db)

	// Ejecutar el test con un ID que no existe
	_, err := repo.GetById(999)

	// Verificar que se devuelve un error de no encontrado
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
}

// Test para GetCantSellersByLocality - Caso exitoso con sellers
func TestLocalitySql_GetCantSellersByLocality_Success(t *testing.T) {
	db := createTestLocalityDB()
	defer db.Close()

	repo := NewLocalitySql(db)
	sellerRepo := NewSellerSql(db)

	// Crear una locality de prueba
	locality := models.Locality{
		LocalitiesAttributes: models.LocalitiesAttributes{
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		},
	}

	createdLocality, err := repo.Create(locality)
	require.NoError(t, err)

	// Crear algunos sellers asociados a la locality
	sellers := []models.Seller{
		{
			SellerAttributes: models.SellerAttributes{
				CId:         "S001",
				CompanyName: "Test Company 1",
				Address:     "123 Test St",
				Telephone:   "123-456-7890",
				LocalityID:  createdLocality.ID,
			},
		},
		{
			SellerAttributes: models.SellerAttributes{
				CId:         "S002",
				CompanyName: "Test Company 2",
				Address:     "456 Test Ave",
				Telephone:   "987-654-3210",
				LocalityID:  createdLocality.ID,
			},
		},
	}

	// Crear los sellers
	for _, seller := range sellers {
		_, err := sellerRepo.Create(seller)
		require.NoError(t, err)
	}

	// Ejecutar el test
	result, err := repo.GetCantSellersByLocality(createdLocality.ID)

	// Verificar resultado
	assert.NoError(t, err)
	assert.Equal(t, createdLocality.ID, result.ID)
	assert.Equal(t, createdLocality.LocalityName, *result.LocalityName)
	assert.Equal(t, "2", *result.SellerCount)
}

// Test para GetCantSellersByLocality - Caso sin sellers
func TestLocalitySql_GetCantSellersByLocality_NoSellers(t *testing.T) {
	db := createTestLocalityDB()
	defer db.Close()

	repo := NewLocalitySql(db)

	// Crear una locality de prueba
	locality := models.Locality{
		LocalitiesAttributes: models.LocalitiesAttributes{
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		},
	}

	createdLocality, err := repo.Create(locality)
	require.NoError(t, err)

	// Ejecutar el test sin crear sellers
	result, err := repo.GetCantSellersByLocality(createdLocality.ID)

	// Verificar resultado
	assert.NoError(t, err)
	assert.Equal(t, createdLocality.ID, result.ID)
	assert.Equal(t, createdLocality.LocalityName, *result.LocalityName)
	assert.Equal(t, "0", *result.SellerCount)
}

// Test para GetCantSellersByLocality - Caso de locality no encontrada
func TestLocalitySql_GetCantSellersByLocality_NotFound(t *testing.T) {
	db := createTestLocalityDB()
	defer db.Close()

	repo := NewLocalitySql(db)

	// Ejecutar el test con un ID que no existe
	_, err := repo.GetCantSellersByLocality(999)

	// Verificar que se devuelve un error de no encontrado
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
}

// Test para verificar integridad referencial - Crear seller con locality inexistente
func TestLocalitySql_ReferentialIntegrity_SellerWithoutLocality(t *testing.T) {
	db := createTestLocalityDB()
	defer db.Close()

	sellerRepo := NewSellerSql(db)

	// Intentar crear un seller con una locality_id que no existe
	seller := models.Seller{
		SellerAttributes: models.SellerAttributes{
			CId:         "S001",
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "123-456-7890",
			LocalityID:  999, // ID que no existe
		},
	}

	// Ejecutar el test
	_, err := sellerRepo.Create(seller)

	// Verificar que se devuelve un error debido a la restricción de clave foránea
	assert.Error(t, err)
}

// Test para verificar que se puede crear un seller con locality existente
func TestLocalitySql_ReferentialIntegrity_SellerWithExistingLocality(t *testing.T) {
	db := createTestLocalityDB()
	defer db.Close()

	repo := NewLocalitySql(db)
	sellerRepo := NewSellerSql(db)

	// Crear una locality
	locality := models.Locality{
		LocalitiesAttributes: models.LocalitiesAttributes{
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		},
	}

	createdLocality, err := repo.Create(locality)
	require.NoError(t, err)

	// Crear un seller con la locality existente
	seller := models.Seller{
		SellerAttributes: models.SellerAttributes{
			CId:         "S001",
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "123-456-7890",
			LocalityID:  createdLocality.ID,
		},
	}

	// Ejecutar el test
	result, err := sellerRepo.Create(seller)

	// Verificar resultado
	assert.NoError(t, err)
	assert.NotZero(t, result.ID)
	assert.Equal(t, seller.LocalityID, result.LocalityID)
}
