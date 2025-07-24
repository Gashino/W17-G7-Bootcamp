package seller

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

// MockSellerService es un mock del servicio de sellers
type MockSellerService struct {
	mock.Mock
}

func (m *MockSellerService) FindAll() (map[int]models.Seller, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Seller), args.Error(1)
}

func (m *MockSellerService) GetById(id int) (models.Seller, error) {
	args := m.Called(id)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerService) Create(seller models.Seller) (models.Seller, error) {
	args := m.Called(seller)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerService) UpdateFields(id int, request models.SellerCreateRequest) (models.Seller, error) {
	args := m.Called(id, request)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerService) DeleteSeller(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
