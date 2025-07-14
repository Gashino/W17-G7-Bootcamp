package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Register txdb driver for testing (only if not already registered)
	user := getEnvOrDefaultPO("DB_USER", "root")
	password := getEnvOrDefaultPO("DB_PASSWORD", "asda1125")
	host := getEnvOrDefaultPO("DB_HOST", "localhost")
	port := getEnvOrDefaultPO("DB_PORT", "3306")
	dbname := getEnvOrDefaultPO("DB_NAME", "db_test")

	cfg := mysql.Config{
		User:                 user,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", host, port),
		DBName:               dbname,
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	// Use a different name to avoid conflict with buyer tests
	txdb.Register("po_txdb", "mysql", cfg.FormatDSN())
}

func getEnvOrDefaultPO(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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

func createTestBuyer(t *testing.T, buyerRepo *BuyerSQL, cardNumberID string) models.Buyer {
	buyer := models.Buyer{
		BuyerAttributes: models.BuyerAttributes{
			CardNumberID: cardNumberID,
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

		purchaseOrder := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "PO-001",
				OrderDate:       "2024-01-15",
				TrackingCode:    "TRK-12345",
				BuyerID:         testBuyer.ID,
				ProductRecordID: 1,
			},
		}

		createdPO, err := repo.Create(purchaseOrder)

		assert.NoError(t, err)
		assert.NotZero(t, createdPO.ID)
		assert.Equal(t, "PO-001", createdPO.OrderNumber)
		assert.Equal(t, "2024-01-15", createdPO.OrderDate)
		assert.Equal(t, "TRK-12345", createdPO.TrackingCode)
		assert.Equal(t, testBuyer.ID, createdPO.BuyerID)
		assert.Equal(t, 1, createdPO.ProductRecordID)
	})

	t.Run("duplicate order_number", func(t *testing.T) {
		// Create a test buyer first
		testBuyer := createTestBuyer(t, buyerRepo, "87654321")

		// First purchase order
		po1 := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "PO-002",
				OrderDate:       "2024-01-15",
				TrackingCode:    "TRK-11111",
				BuyerID:         testBuyer.ID,
				ProductRecordID: 1,
			},
		}
		_, err := repo.Create(po1)
		require.NoError(t, err)

		// Second purchase order with same order_number
		po2 := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "PO-002", // Duplicate
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
		purchaseOrder := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "PO-003",
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

		// Create first purchase order
		po1 := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "PO-004",
				OrderDate:       "2024-01-15",
				TrackingCode:    "TRK-44444",
				BuyerID:         testBuyer.ID,
				ProductRecordID: 1,
			},
		}
		createdPO1, err := repo.Create(po1)
		require.NoError(t, err)

		// Create second purchase order for same buyer
		po2 := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "PO-005",
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

		purchaseOrder := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "PO-006",
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

		purchaseOrder := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "PO-007",
				OrderDate:       "2024-12-31", // Different date format
				TrackingCode:    "TRK-77777",
				BuyerID:         testBuyer.ID,
				ProductRecordID: 1,
			},
		}

		createdPO, err := repo.Create(purchaseOrder)

		assert.NoError(t, err)
		assert.Equal(t, "2024-12-31", createdPO.OrderDate)
	})
}
