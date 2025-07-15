package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Estructura para leer la configuración de la base de datos desde config.yml
type PODBConfig struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

// Cargar la configuración desde config.yml
func loadPOConfig() PODBConfig {
	var config PODBConfig

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

func init() {
	// Leer la configuración desde config.yml
	config := loadPOConfig()

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

	// Use a different name to avoid conflict with buyer tests
	txdb.Register("po_txdb", "mysql", cfg.FormatDSN())
}

func setupPurchaseOrderTestDB(t *testing.T) (*PurchaseOrderSQL, *BuyerSQL) {
	db, err := sql.Open("po_txdb", "purchase_order_test")
	require.NoError(t, err)

	// Create buyers table (needed for foreign key constraint)
	createBuyersTableQuery := `
		CREATE TABLE IF NOT EXISTS buyers (
			id INT AUTO_INCREMENT PRIMARY KEY,
			card_number_id VARCHAR(8) UNIQUE NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`
	_, err = db.Exec(createBuyersTableQuery)
	require.NoError(t, err)

	// Create purchase_orders table
	createPurchaseOrdersTableQuery := `
		CREATE TABLE IF NOT EXISTS purchase_orders (
			id INT AUTO_INCREMENT PRIMARY KEY,
			order_number VARCHAR(50) UNIQUE NOT NULL,
			order_date DATE NOT NULL,
			tracking_code VARCHAR(100) NOT NULL,
			buyer_id INT NOT NULL,
			product_record_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (buyer_id) REFERENCES buyers(id)
		)
	`
	_, err = db.Exec(createPurchaseOrdersTableQuery)
	require.NoError(t, err)

	poRepo := NewPurchaseOrderSQL(db)
	buyerRepo := NewBuyerSQL(db)
	return poRepo, buyerRepo
}

// Variables para generar IDs únicos para cada test
var poCardCounter int = 0
var orderCounter int = 0

// Función para generar un número de tarjeta único para purchase orders
func generatePOUniqueCardID() string {
	poCardCounter++
	return fmt.Sprintf("POCARD-%d-%d", time.Now().UnixNano(), poCardCounter)
}

// Función para generar un número de orden único
func generateUniqueOrderNumber() string {
	orderCounter++
	return fmt.Sprintf("PO-%d-%d", time.Now().UnixNano(), orderCounter)
}

func createTestBuyer(t *testing.T, buyerRepo *BuyerSQL, cardNumberIDSuffix string) models.Buyer {
	// Generar un ID único para cada test
	uniqueCardID := generatePOUniqueCardID() + "-" + cardNumberIDSuffix
	
	buyer := models.Buyer{
		BuyerAttributes: models.BuyerAttributes{
			CardNumberID: uniqueCardID,
			FirstName:    "Test",
			LastName:     "Buyer",
		},
	}

	createdBuyer, err := buyerRepo.Create(buyer)
	require.NoError(t, err)
	return createdBuyer
}

func TestPurchaseOrderSQL_Create(t *testing.T) {
	repo, buyerRepo := setupPurchaseOrderTestDB(t)

	t.Run("successful creation", func(t *testing.T) {
		// Create a test buyer first
		testBuyer := createTestBuyer(t, buyerRepo, "12345678")

		// Generar un número de orden único para este test
		uniqueOrderNumber := generateUniqueOrderNumber()

		purchaseOrder := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     uniqueOrderNumber,
				OrderDate:       "2024-01-15",
				TrackingCode:    "TRK-12345",
				BuyerID:         testBuyer.ID,
				ProductRecordID: 1,
			},
		}

		createdPO, err := repo.Create(purchaseOrder)

		assert.NoError(t, err)
		assert.NotZero(t, createdPO.ID)
		assert.Equal(t, uniqueOrderNumber, createdPO.OrderNumber)
		assert.Equal(t, "2024-01-15", createdPO.OrderDate)
		assert.Equal(t, "TRK-12345", createdPO.TrackingCode)
		assert.Equal(t, testBuyer.ID, createdPO.BuyerID)
		assert.Equal(t, 1, createdPO.ProductRecordID)
	})

	t.Run("duplicate order_number", func(t *testing.T) {
		// Create a test buyer first
		testBuyer := createTestBuyer(t, buyerRepo, "87654321")

		// Generar un número de orden único para el primer purchase order
		uniqueOrderNumber1 := generateUniqueOrderNumber()

		// First purchase order
		po1 := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     uniqueOrderNumber1,
				OrderDate:       "2024-01-15",
				TrackingCode:    "TRK-11111",
				BuyerID:         testBuyer.ID,
				ProductRecordID: 1,
			},
		}
		_, err := repo.Create(po1)
		require.NoError(t, err)

		// Second purchase order with same order_number (para testear el caso de duplicado)
		po2 := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     uniqueOrderNumber1, // Usar el mismo número que po1 para testear duplicados
				OrderDate:       "2024-01-16",
				TrackingCode:    "TRK-22222",
				BuyerID:         testBuyer.ID,
				ProductRecordID: 2,
			},
		}

		_, err = repo.Create(po2)
		assert.Error(t, err)

		serviceError, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, serviceError.ResponseCode)
		assert.Contains(t, serviceError.InternalError.Error(), "order_number already exists")
	})

	t.Run("invalid buyer_id (foreign key constraint)", func(t *testing.T) {
		// Generar un número de orden único para este test
		uniqueOrderNumber := generateUniqueOrderNumber()

		purchaseOrder := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     uniqueOrderNumber,
				OrderDate:       "2024-01-15",
				TrackingCode:    "TRK-33333",
				BuyerID:         99999, // Non-existent buyer
				ProductRecordID: 1,
			},
		}

		_, err := repo.Create(purchaseOrder)
		assert.Error(t, err)

		serviceError, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, serviceError.ResponseCode)
		assert.Contains(t, serviceError.InternalError.Error(), "buyer_id does not exist")
	})

	t.Run("multiple purchase orders for same buyer", func(t *testing.T) {
		// Create a test buyer first
		testBuyer := createTestBuyer(t, buyerRepo, "11111111")

		// Generar números de orden únicos para este test
		uniqueOrderNumber1 := generateUniqueOrderNumber()

		// Create first purchase order
		po1 := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     uniqueOrderNumber1,
				OrderDate:       "2024-01-15",
				TrackingCode:    "TRK-44444",
				BuyerID:         testBuyer.ID,
				ProductRecordID: 1,
			},
		}
		createdPO1, err := repo.Create(po1)
		require.NoError(t, err)

		// Generar otro número de orden único
		uniqueOrderNumber2 := generateUniqueOrderNumber()

		// Create second purchase order for same buyer
		po2 := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     uniqueOrderNumber2,
				OrderDate:       "2024-01-16",
				TrackingCode:    "TRK-55555",
				BuyerID:         testBuyer.ID,
				ProductRecordID: 2,
			},
		}
		createdPO2, err := repo.Create(po2)

		assert.NoError(t, err)
		assert.NotEqual(t, createdPO1.ID, createdPO2.ID)
		assert.Equal(t, testBuyer.ID, createdPO1.BuyerID)
		assert.Equal(t, testBuyer.ID, createdPO2.BuyerID)
		assert.NotEqual(t, createdPO1.OrderNumber, createdPO2.OrderNumber)
	})

	t.Run("edge case - empty strings should not cause issues", func(t *testing.T) {
		// Create a test buyer first
		testBuyer := createTestBuyer(t, buyerRepo, "22222222")

		// Generar un número de orden único para este test
		uniqueOrderNumber := generateUniqueOrderNumber()

		purchaseOrder := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     uniqueOrderNumber,
				OrderDate:       "2024-01-15",
				TrackingCode:    "", // Empty tracking code should be allowed
				BuyerID:         testBuyer.ID,
				ProductRecordID: 1,
			},
		}

		createdPO, err := repo.Create(purchaseOrder)

		assert.NoError(t, err)
		assert.Equal(t, "", createdPO.TrackingCode)
	})

	t.Run("different date formats", func(t *testing.T) {
		// Create a test buyer first
		testBuyer := createTestBuyer(t, buyerRepo, "33333333")

		// Generar un número de orden único para este test
		uniqueOrderNumber := generateUniqueOrderNumber()

		purchaseOrder := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     uniqueOrderNumber,
				OrderDate:       "2024-12-31", // Different date format
				TrackingCode:    "TRK-77777",
				BuyerID:         testBuyer.ID,
				ProductRecordID: 1,
			},
		}

		createdPO, err := repo.Create(purchaseOrder)

		assert.NoError(t, err)
		assert.Equal(t, "2024-12-31", createdPO.OrderDate)
		assert.Equal(t, uniqueOrderNumber, createdPO.OrderNumber)
	})
}
