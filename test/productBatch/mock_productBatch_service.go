package productbatch

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductBatchService struct {
	mock.Mock
}

func (m *MockProductBatchService) PostProductBatch(batch models.ProductBatch) (models.ProductBatch, error) {
	args := m.Called(batch)
	return args.Get(0).(models.ProductBatch), args.Error(1)
}
