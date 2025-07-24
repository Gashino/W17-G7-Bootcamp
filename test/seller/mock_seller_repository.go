package seller

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

// MockSellerRepository es un mock del repositorio de sellers
type MockSellerRepository struct {
	mock.Mock
}

func (m *MockSellerRepository) FindAll() (map[int]models.Seller, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Seller), args.Error(1)
}

func (m *MockSellerRepository) GetById(id int) (models.Seller, error) {
	args := m.Called(id)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerRepository) Create(seller models.Seller) (models.Seller, error) {
	args := m.Called(seller)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerRepository) UpdateFields(id int, data models.SellerCreateRequest) (models.Seller, error) {
	args := m.Called(id, data)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerRepository) DeleteSeller(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
