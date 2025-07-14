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
)

var (
	cardIDCounter int64
	cardIDMutex   sync.Mutex
)

func init() {
	// Register txdb driver for testing
	user := getEnvOrDefault("DB_USER", "root")
	password := getEnvOrDefault("DB_PASSWORD", "asda1125")
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "3306")
	dbname := getEnvOrDefault("DB_NAME", "db_test")

	cfg := mysql.Config{
		User:                 user,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", host, port),
		DBName:               dbname,
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// generateUniqueCardID generates a unique 8-digit card ID using timestamp and counter
func generateUniqueCardID() string {
	cardIDMutex.Lock()
	defer cardIDMutex.Unlock()

	cardIDCounter++
	// Add a small sleep to ensure different timestamps
	time.Sleep(1 * time.Millisecond)

	// Combine timestamp and counter for guaranteed uniqueness
	unique := (time.Now().UnixNano() % 100000000) + cardIDCounter
	return fmt.Sprintf("%08d", unique%100000000)
}

func setupBuyerTestDB(t *testing.T) (*BuyerSQL, *sql.DB) {
	db, err := sql.Open("txdb", "buyer_test")
	require.NoError(t, err)

	// Create buyers table if it doesn't exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS buyers (
			id INT AUTO_INCREMENT PRIMARY KEY,
			card_number_id VARCHAR(8) UNIQUE NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`
	_, err = db.Exec(createTableQuery)
	require.NoError(t, err)

	return NewBuyerSQL(db), db
}

func TestBuyerSQL_Create(t *testing.T) {
	repo, _ := setupBuyerTestDB(t)

	t.Run("successful creation", func(t *testing.T) {
		buyer := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: generateUniqueCardID(),
				FirstName:    "John",
				LastName:     "Doe",
			},
		}

		createdBuyer, err := repo.Create(buyer)

		assert.NoError(t, err)
		assert.NotZero(t, createdBuyer.ID)
		assert.Equal(t, buyer.CardNumberID, createdBuyer.CardNumberID)
		assert.Equal(t, "John", createdBuyer.FirstName)
		assert.Equal(t, "Doe", createdBuyer.LastName)
	})

	t.Run("duplicate card_number_id", func(t *testing.T) {
		cardID := generateUniqueCardID()

		// First buyer
		buyer1 := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: cardID,
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		}
		_, err := repo.Create(buyer1)
		require.NoError(t, err)

		// Second buyer with same card_number_id
		buyer2 := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: cardID, // Duplicate
				FirstName:    "Alice",
				LastName:     "Johnson",
			},
		}

		_, err = repo.Create(buyer2)
		assert.Error(t, err)

		serviceError, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, serviceError.ResponseCode)
	})
}

func TestBuyerSQL_GetByID(t *testing.T) {
	repo, _ := setupBuyerTestDB(t)

	t.Run("existing buyer", func(t *testing.T) {
		// Create a buyer first
		buyer := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: generateUniqueCardID(),
				FirstName:    "Test",
				LastName:     "User",
			},
		}
		createdBuyer, err := repo.Create(buyer)
		require.NoError(t, err)

		// Retrieve the buyer
		retrievedBuyer, err := repo.GetByID(createdBuyer.ID)

		assert.NoError(t, err)
		assert.Equal(t, createdBuyer.ID, retrievedBuyer.ID)
		assert.Equal(t, buyer.CardNumberID, retrievedBuyer.CardNumberID)
		assert.Equal(t, "Test", retrievedBuyer.FirstName)
		assert.Equal(t, "User", retrievedBuyer.LastName)
	})

	t.Run("non-existing buyer", func(t *testing.T) {
		_, err := repo.GetByID(99999)

		assert.Error(t, err)
		serviceError, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, serviceError.ResponseCode)
	})
}

func TestBuyerSQL_GetAll(t *testing.T) {
	repo, db := setupBuyerTestDB(t)

	t.Run("empty table", func(t *testing.T) {
		// Clean table for this test
		_, err := db.Exec("DELETE FROM buyers")
		require.NoError(t, err)

		buyers, err := repo.GetAll()

		assert.NoError(t, err)
		assert.Empty(t, buyers)
	})

	t.Run("with buyers", func(t *testing.T) {
		// Clean table for this test
		_, err := db.Exec("DELETE FROM buyers")
		require.NoError(t, err)

		// Create multiple buyers
		buyer1 := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: generateUniqueCardID(),
				FirstName:    "Alice",
				LastName:     "Wonder",
			},
		}
		buyer2 := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: generateUniqueCardID(),
				FirstName:    "Bob",
				LastName:     "Builder",
			},
		}

		createdBuyer1, err := repo.Create(buyer1)
		require.NoError(t, err)
		createdBuyer2, err := repo.Create(buyer2)
		require.NoError(t, err)

		// Retrieve all buyers
		buyers, err := repo.GetAll()

		assert.NoError(t, err)
		assert.Len(t, buyers, 2)
		assert.Contains(t, buyers, createdBuyer1.ID)
		assert.Contains(t, buyers, createdBuyer2.ID)
		assert.Equal(t, "Alice", buyers[createdBuyer1.ID].FirstName)
		assert.Equal(t, "Bob", buyers[createdBuyer2.ID].FirstName)
	})
}

func TestBuyerSQL_Update(t *testing.T) {
	repo, _ := setupBuyerTestDB(t)

	t.Run("successful update", func(t *testing.T) {
		// Create a buyer first
		buyer := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: generateUniqueCardID(),
				FirstName:    "Original",
				LastName:     "Name",
			},
		}
		createdBuyer, err := repo.Create(buyer)
		require.NoError(t, err)

		// Update the buyer
		updateBuyer := models.Buyer{
			ID: createdBuyer.ID,
			BuyerAttributes: models.BuyerAttributes{
				FirstName: "Updated",
				LastName:  "NameChanged",
			},
		}

		updatedBuyer, err := repo.Update(updateBuyer)

		assert.NoError(t, err)
		assert.Equal(t, createdBuyer.ID, updatedBuyer.ID)
		assert.Equal(t, buyer.CardNumberID, updatedBuyer.CardNumberID) // Should remain unchanged
		assert.Equal(t, "Updated", updatedBuyer.FirstName)
		assert.Equal(t, "NameChanged", updatedBuyer.LastName)
	})

	t.Run("partial update", func(t *testing.T) {
		// Create a buyer first
		buyer := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: generateUniqueCardID(),
				FirstName:    "Partial",
				LastName:     "Update",
			},
		}
		createdBuyer, err := repo.Create(buyer)
		require.NoError(t, err)

		// Update only first name
		updateBuyer := models.Buyer{
			ID: createdBuyer.ID,
			BuyerAttributes: models.BuyerAttributes{
				FirstName: "OnlyFirst",
				// LastName and CardNumberID not specified (empty)
			},
		}

		updatedBuyer, err := repo.Update(updateBuyer)

		assert.NoError(t, err)
		assert.Equal(t, "OnlyFirst", updatedBuyer.FirstName)
		assert.Equal(t, "Update", updatedBuyer.LastName)               // Should remain unchanged
		assert.Equal(t, buyer.CardNumberID, updatedBuyer.CardNumberID) // Should remain unchanged
	})

	t.Run("non-existing buyer", func(t *testing.T) {
		updateBuyer := models.Buyer{
			ID: 99999,
			BuyerAttributes: models.BuyerAttributes{
				FirstName: "NonExisting",
			},
		}

		_, err := repo.Update(updateBuyer)

		assert.Error(t, err)
		serviceError, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, serviceError.ResponseCode)
	})

	t.Run("duplicate card_number_id", func(t *testing.T) {
		// Create two buyers
		cardID1 := generateUniqueCardID()
		cardID2 := generateUniqueCardID()

		buyer1 := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: cardID1,
				FirstName:    "First",
				LastName:     "Buyer",
			},
		}
		buyer2 := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: cardID2,
				FirstName:    "Second",
				LastName:     "Buyer",
			},
		}

		createdBuyer1, err := repo.Create(buyer1)
		require.NoError(t, err)
		createdBuyer2, err := repo.Create(buyer2)
		require.NoError(t, err)

		// Try to update buyer2 with buyer1's card_number_id
		updateBuyer := models.Buyer{
			ID: createdBuyer2.ID,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: createdBuyer1.CardNumberID, // Use existing card_number_id
			},
		}

		_, err = repo.Update(updateBuyer)

		assert.Error(t, err)
		serviceError, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, serviceError.ResponseCode)
	})
}

func TestBuyerSQL_Delete(t *testing.T) {
	repo, _ := setupBuyerTestDB(t)

	t.Run("successful deletion", func(t *testing.T) {
		// Create a buyer first
		buyer := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: generateUniqueCardID(),
				FirstName:    "ToDelete",
				LastName:     "User",
			},
		}
		createdBuyer, err := repo.Create(buyer)
		require.NoError(t, err)

		// Delete the buyer
		err = repo.Delete(createdBuyer.ID)
		assert.NoError(t, err)

		// Verify buyer is deleted
		_, err = repo.GetByID(createdBuyer.ID)
		assert.Error(t, err)
		serviceError, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, serviceError.ResponseCode)
	})

	t.Run("non-existing buyer", func(t *testing.T) {
		err := repo.Delete(99999)

		assert.Error(t, err)
		serviceError, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, serviceError.ResponseCode)
	})
}

func TestBuyerSQL_GetPurchaseOrdersReport(t *testing.T) {
	repo, _ := setupBuyerTestDB(t)

	t.Run("all buyers report - no purchase orders", func(t *testing.T) {
		// Create a buyer without any purchase orders
		buyer := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: generateUniqueCardID(),
				FirstName:    "Report",
				LastName:     "Test",
			},
		}
		createdBuyer, err := repo.Create(buyer)
		require.NoError(t, err)

		// Get all buyers report
		reports, err := repo.GetPurchaseOrdersReport(nil)

		assert.NoError(t, err)
		assert.NotEmpty(t, reports)

		// Find our buyer in the reports
		var foundBuyer *models.BuyerPurchaseOrderReport
		for _, report := range reports {
			if report.ID == createdBuyer.ID {
				foundBuyer = &report
				break
			}
		}

		assert.NotNil(t, foundBuyer)
		assert.Equal(t, createdBuyer.ID, foundBuyer.ID)
		assert.Equal(t, buyer.CardNumberID, foundBuyer.CardNumberID)
		assert.Equal(t, "Report", foundBuyer.FirstName)
		assert.Equal(t, "Test", foundBuyer.LastName)
		assert.Equal(t, 0, foundBuyer.PurchaseOrdersCount)
	})

	t.Run("specific buyer report - no purchase orders", func(t *testing.T) {
		// Create a buyer without any purchase orders
		buyer := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: generateUniqueCardID(),
				FirstName:    "Specific",
				LastName:     "Test",
			},
		}
		createdBuyer, err := repo.Create(buyer)
		require.NoError(t, err)

		// Get specific buyer report
		reports, err := repo.GetPurchaseOrdersReport(&createdBuyer.ID)

		assert.NoError(t, err)
		assert.Len(t, reports, 1)
		assert.Equal(t, createdBuyer.ID, reports[0].ID)
		assert.Equal(t, buyer.CardNumberID, reports[0].CardNumberID)
		assert.Equal(t, "Specific", reports[0].FirstName)
		assert.Equal(t, "Test", reports[0].LastName)
		assert.Equal(t, 0, reports[0].PurchaseOrdersCount)
	})

	t.Run("specific buyer report - non-existing buyer", func(t *testing.T) {
		nonExistingID := 99999

		reports, err := repo.GetPurchaseOrdersReport(&nonExistingID)

		assert.Error(t, err)
		assert.Empty(t, reports)
		serviceError, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, serviceError.ResponseCode)
	})

	// Note: To fully test with purchase orders, we would need to create
	// the purchase_orders table and insert test data. This would require
	// extending the test setup to include purchase order creation.
}
