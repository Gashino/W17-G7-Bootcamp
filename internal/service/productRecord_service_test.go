package service

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/productRecord"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Helper functions for creating pointers
func timePtr(t time.Time) *time.Time { return &t }

// Function to create a test product record
func createTestProductRecord() models.ProductRecord {
	now := time.Now()
	return models.ProductRecord{
		ID:             1,
		LastUpdateDate: timePtr(now),
		PurchasePrice:  float64Ptr(100.50),
		SalePrice:      float64Ptr(150.75),
		ProductId:      models.IntPtr(1),
	}
}

func TestProductRecordService_Create(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		// Arrange
		mockRepo := new(productRecord.MockProductRecordRepository)
		productRecordService := NewProductRecordDefault(mockRepo)

		newRecord := createTestProductRecord()
		newRecord.ID = 0 // ID should be assigned by the repository

		expectedRecord := createTestProductRecord() // ID = 1

		mockRepo.On("Insert", newRecord).Return(&expectedRecord, nil)

		// Act
		record, err := productRecordService.Create(newRecord)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &expectedRecord, record)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_with_invalid_product_id", func(t *testing.T) {
		// Arrange
		mockRepo := new(productRecord.MockProductRecordRepository)
		productRecordService := NewProductRecordDefault(mockRepo)

		newRecord := createTestProductRecord()
		newRecord.ID = 0
		*newRecord.ProductId = 999 // Non-existent product ID

		errorResponse := pkg.ServiceErrors[pkg.ErrNotFound]
		errorResponse.Message = "invalid product_id"
		mockRepo.On("Insert", newRecord).Return(nil, errorResponse)

		// Act
		record, err := productRecordService.Create(newRecord)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Equal(t, errorResponse, err)
		mockRepo.AssertCalled(t, "Insert", newRecord)
	})

	t.Run("create_with_internal_error", func(t *testing.T) {
		// Arrange
		mockRepo := new(productRecord.MockProductRecordRepository)
		productRecordService := NewProductRecordDefault(mockRepo)

		newRecord := createTestProductRecord()
		newRecord.ID = 0

		internalError := pkg.ServiceErrors[pkg.ErrInternalServer]
		mockRepo.On("Insert", newRecord).Return(nil, internalError)

		// Act
		record, err := productRecordService.Create(newRecord)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Equal(t, internalError, err)
		mockRepo.AssertCalled(t, "Insert", newRecord)
	})

	t.Run("create_with_database_error", func(t *testing.T) {
		// Arrange
		mockRepo := new(productRecord.MockProductRecordRepository)
		productRecordService := NewProductRecordDefault(mockRepo)

		newRecord := createTestProductRecord()
		newRecord.ID = 0

		databaseError := errors.New("database connection error")
		mockRepo.On("Insert", newRecord).Return(nil, databaseError)

		// Act
		record, err := productRecordService.Create(newRecord)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Equal(t, databaseError, err)
		mockRepo.AssertCalled(t, "Insert", newRecord)
	})
}
