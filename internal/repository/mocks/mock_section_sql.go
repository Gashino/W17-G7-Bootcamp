package mocks

import (
	"app/pkg/models"
	"github.com/stretchr/testify/mock"
)

type MockSectionRepository struct {
	mock.Mock
}

func (m *MockSectionRepository) GetAll() (map[int]models.Section, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Section), args.Error(1)
}

func (m *MockSectionRepository) GetByID(id int) (models.Section, error) {
	args := m.Called(id)
	return args.Get(0).(models.Section), args.Error(1)
}

func (m *MockSectionRepository) Create(section models.Section) (models.Section, error) {
	args := m.Called(section)
	return args.Get(0).(models.Section), args.Error(1)
}

func (m *MockSectionRepository) Update(id int, section models.Section) (models.Section, error) {
	args := m.Called(id, section)
	return args.Get(0).(models.Section), args.Error(1)
}

func (m *MockSectionRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
