package locality

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

// MockLocalityService is a mock of locality service
type MockLocalityService struct {
	mock.Mock
}

func (m *MockLocalityService) Create(locality models.Locality) (models.Locality, error) {
	args := m.Called(locality)
	return args.Get(0).(models.Locality), args.Error(1)
}

func (m *MockLocalityService) GetById(id int) (models.Locality, error) {
	args := m.Called(id)
	return args.Get(0).(models.Locality), args.Error(1)
}

func (m *MockLocalityService) GetCantSellersByLocality(id int) (models.LocalityBySellerResponse, error) {
	args := m.Called(id)
	return args.Get(0).(models.LocalityBySellerResponse), args.Error(1)
}
