package productRecord

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductRecordService struct {
	mock.Mock
}

// Create implements service.IProductRecordService.
func (m *MockProductRecordService) Create(productRecord models.ProductRecord) (*models.ProductRecord, error) {
	args := m.Called(productRecord)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductRecord), args.Error(1)
}
