package locality

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

// MockLocalityRepository is a mock implementation of LocalityRepository
type MockLocalityRepository struct {
	mock.Mock
}

// GetById mocks the GetById method
func (m *MockLocalityRepository) GetById(id int) (models.Locality, error) {
	args := m.Called(id)
	return args.Get(0).(models.Locality), args.Error(1)
}

// Create mocks the Create method
func (m *MockLocalityRepository) Create(locality models.Locality) (models.Locality, error) {
	args := m.Called(locality)
	return args.Get(0).(models.Locality), args.Error(1)
}

// GetCantSellersByLocality mocks the GetCantSellersByLocality method
func (m *MockLocalityRepository) GetCantSellersByLocality(id int) (models.LocalityBySellerResponse, error) {
	args := m.Called(id)
	return args.Get(0).(models.LocalityBySellerResponse), args.Error(1)
}
