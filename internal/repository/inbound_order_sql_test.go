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

// Variables para controlar el registro del driver txdb para inbound order tests
var (
	inboundOrderTxdbRegistered bool = false
	inboundOrderTxdbMutex      sync.Mutex
)

// Estructura para leer la configuración de la base de datos desde config.yml
type InboundOrderConfig struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

// Configurar la base de datos para pruebas de inbound order
func setupInboundOrderTxDB() string {
	// Usar un mutex para evitar condiciones de carrera al registrar el driver
	inboundOrderTxdbMutex.Lock()
	defer inboundOrderTxdbMutex.Unlock()

	// Evitar registrar el driver más de una vez
	if !inboundOrderTxdbRegistered {
		// Leer la configuración desde config.yml
		config := loadInboundOrderConfig()

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
			if driver == "inbound_order_txdb" {
				// El driver ya está registrado, no necesitamos registrarlo de nuevo
				inboundOrderTxdbRegistered = true
				return "inbound_order_txdb"
			}
		}

		// Registrar el driver si no está registrado
		txdb.Register("inbound_order_txdb", "mysql", cfg.FormatDSN())
		inboundOrderTxdbRegistered = true
	}

	return "inbound_order_txdb"
}

// Cargar la configuración desde config.yml para inbound order tests
func loadInboundOrderConfig() InboundOrderConfig {
	var config InboundOrderConfig

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

// Configurar la base de datos de prueba para cada test de inbound order
func setupInboundOrderTestDB(t *testing.T) (*InboundOrderSQL, *sql.DB) {
	// Usar un nombre único para cada conexión de test
	driver := setupInboundOrderTxDB()
	db, err := sql.Open(driver, fmt.Sprintf("inbound_order_test_%s", t.Name()))
	require.NoError(t, err)

	// Verificar que la conexión funciona
	err = db.Ping()
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	return &InboundOrderSQL{db: db}, db
}

// Crear una inbound order completa para pruebas
func createTestInboundOrder() models.InboundOrder {
	// Usar timestamp para generar un order number único
	timestamp := time.Now().UnixNano()
	orderNumber := fmt.Sprintf("ORD-TEST-%d", timestamp)

	return models.InboundOrder{
		OrderDate:      time.Now(),
		OrderNumber:    orderNumber,
		EmployeeID:     1, // Usar un employee_id válido según el script de base de datos
		ProductBatchID: 1, // Usar un product_batch_id válido según el script de base de datos
		WarehouseID:    1, // Usar un warehouse_id válido según el script de base de datos
	}
}

func TestInboundOrderSQL_Create(t *testing.T) {
	t.Run("create new inbound order successfully", func(t *testing.T) {
		repo, db := setupInboundOrderTestDB(t)
		defer db.Close()

		// Arrange
		newInboundOrder := createTestInboundOrder()

		// Act
		createdInboundOrder, err := repo.Create(newInboundOrder)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, createdInboundOrder)
		assert.Equal(t, newInboundOrder.OrderNumber, createdInboundOrder.OrderNumber)
		assert.Equal(t, newInboundOrder.EmployeeID, createdInboundOrder.EmployeeID)
		assert.Equal(t, newInboundOrder.ProductBatchID, createdInboundOrder.ProductBatchID)
		assert.Equal(t, newInboundOrder.WarehouseID, createdInboundOrder.WarehouseID)
		// La fecha debería ser similar (permitir algunos segundos de diferencia)
		assert.WithinDuration(t, newInboundOrder.OrderDate, createdInboundOrder.OrderDate, 5*time.Second)
	})

	t.Run("create inbound order with invalid employee id", func(t *testing.T) {
		repo, db := setupInboundOrderTestDB(t)
		defer db.Close()

		// Arrange - Crear una inbound order con un employee_id inválido
		invalidInboundOrder := createTestInboundOrder()
		invalidInboundOrder.EmployeeID = 999 // ID que no existe

		// Act
		createdInboundOrder, err := repo.Create(invalidInboundOrder)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.InboundOrder{}, createdInboundOrder)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, err.(pkg.ServiceError).Code)
	})

	t.Run("create inbound order with invalid product batch id", func(t *testing.T) {
		repo, db := setupInboundOrderTestDB(t)
		defer db.Close()

		// Arrange - Crear una inbound order con un product_batch_id inválido
		invalidInboundOrder := createTestInboundOrder()
		invalidInboundOrder.ProductBatchID = 999 // ID que no existe

		// Act
		createdInboundOrder, err := repo.Create(invalidInboundOrder)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.InboundOrder{}, createdInboundOrder)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, err.(pkg.ServiceError).Code)
	})

	t.Run("create inbound order with invalid warehouse id", func(t *testing.T) {
		repo, db := setupInboundOrderTestDB(t)
		defer db.Close()

		// Arrange - Crear una inbound order con un warehouse_id inválido
		invalidInboundOrder := createTestInboundOrder()
		invalidInboundOrder.WarehouseID = 999 // ID que no existe

		// Act
		createdInboundOrder, err := repo.Create(invalidInboundOrder)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.InboundOrder{}, createdInboundOrder)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, err.(pkg.ServiceError).Code)
	})

	t.Run("create inbound order with empty order number", func(t *testing.T) {
		repo, db := setupInboundOrderTestDB(t)
		defer db.Close()

		// Arrange - Crear una inbound order con order_number vacío
		invalidInboundOrder := createTestInboundOrder()
		invalidInboundOrder.OrderNumber = "" // Order number vacío

		// Act
		createdInboundOrder, err := repo.Create(invalidInboundOrder)

		// Assert
		// La base de datos permite order_number vacío (NOT NULL pero acepta strings vacíos)
		// Por lo tanto, este test verifica que se puede crear correctamente
		assert.NoError(t, err)
		assert.NotEmpty(t, createdInboundOrder)
		assert.Equal(t, "", createdInboundOrder.OrderNumber)
	})
}
