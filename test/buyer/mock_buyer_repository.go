package buyer

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockBuyerRepository struct {
	mock.Mock
}

func (m *MockBuyerRepository) GetAll() (map[int]models.Buyer, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) GetByID(id int) (models.Buyer, error) {
	args := m.Called(id)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) Create(buyer models.Buyer) (models.Buyer, error) {
	args := m.Called(buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) Update(buyer models.Buyer) (models.Buyer, error) {
	args := m.Called(buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBuyerRepository) GetPurchaseOrdersReport(buyerID *int) ([]models.BuyerPurchaseOrderReport, error) {
	args := m.Called(buyerID)
	return args.Get(0).([]models.BuyerPurchaseOrderReport), args.Error(1)
}
