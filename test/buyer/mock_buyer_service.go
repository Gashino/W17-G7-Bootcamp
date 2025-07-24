package buyer

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockBuyerService struct {
	mock.Mock
}

func (m *MockBuyerService) GetAll() (map[int]models.Buyer, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Buyer), args.Error(1)
}

func (m *MockBuyerService) GetByID(id int) (models.Buyer, error) {
	args := m.Called(id)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerService) Create(buyer models.Buyer) (models.Buyer, error) {
	args := m.Called(buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerService) Update(id int, buyer models.Buyer) (models.Buyer, error) {
	args := m.Called(id, buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBuyerService) GetPurchaseOrdersReport(buyerID *int) ([]models.BuyerPurchaseOrderReport, error) {
	args := m.Called(buyerID)
	return args.Get(0).([]models.BuyerPurchaseOrderReport), args.Error(1)
}
