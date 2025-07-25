package productRecord

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductRecordRepository struct {
	mock.Mock
}

// Insert implements repository.ProductRecordRepository.
func (m *MockProductRecordRepository) Insert(record models.ProductRecord) (*models.ProductRecord, error) {
	args := m.Called(record)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductRecord), args.Error(1)
}
