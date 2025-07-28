package productbatch

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductBatchRepository struct {
	mock.Mock
}

func (m *MockProductBatchRepository) InsertProductBatch(pb models.ProductBatch) (models.ProductBatch, error) {
	args := m.Called(pb)
	return args.Get(0).(models.ProductBatch), args.Error(1)
}
