package service

import (
	"app/pkg"
	"app/pkg/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSectionRepository es un mock del repositorio de secciones
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

func createTestSection(id int) models.Section {
	return models.Section{
		ID: id,
		SectionAttributes: models.SectionAttributes{
			SectionNumber:      id * 100,
			CurrentTemperature: 20.5,
			MinimumTemperature: 15.0,
			CurrentCapacity:    50,
			MinimumCapacity:    10,
			MaximumCapacity:    100,
			WarehouseID:        1,
			ProductTypeID:      1,
		},
	}
}

func TestSectionService_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		expectedSections := map[int]models.Section{
			1: createTestSection(1),
			2: createTestSection(2),
		}

		repo := new(MockSectionRepository)
		repo.On("GetAll").Return(expectedSections, nil)

		service := NewSectionDefault(repo)

		// Act
		sections, err := service.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, sections, 2)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database error")

		repo := new(MockSectionRepository)
		repo.On("GetAll").Return(map[int]models.Section{}, expectedErr)

		service := NewSectionDefault(repo)

		// Act
		sections, err := service.GetAll()

		// Assert
		assert.Error(t, err)
		assert.Empty(t, sections)
		repo.AssertExpectations(t)
	})
}

func TestSectionService_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		expectedSection := createTestSection(1)

		repo := new(MockSectionRepository)
		repo.On("GetByID", 1).Return(expectedSection, nil)

		service := NewSectionDefault(repo)

		// Act
		section, err := service.GetByID(1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedSection, section)
		repo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Arrange
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = errors.New("section not found")

		repo := new(MockSectionRepository)
		repo.On("GetByID", 999).Return(models.Section{}, svcErr)

		service := NewSectionDefault(repo)

		// Act
		section, err := service.GetByID(999)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, section.ID)
		repo.AssertExpectations(t)
	})
}

func TestSectionService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		newSection := createTestSection(0) // ID será asignado por el repositorio
		expectedSection := createTestSection(1)

		repo := new(MockSectionRepository)
		repo.On("Create", newSection).Return(expectedSection, nil)

		service := NewSectionDefault(repo)

		// Act
		createdSection, err := service.Create(newSection)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedSection, createdSection)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		newSection := createTestSection(0)
		expectedErr := errors.New("create error")

		repo := new(MockSectionRepository)
		repo.On("Create", newSection).Return(models.Section{}, expectedErr)

		service := NewSectionDefault(repo)

		// Act
		section, err := service.Create(newSection)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, section.ID)
		repo.AssertExpectations(t)
	})
}

func TestSectionService_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		updateData := models.Section{
			ID: 1,
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      100,
				CurrentTemperature: 25.0,
			},
		}
		expectedSection := createTestSection(1)
		expectedSection.SectionNumber = 100
		expectedSection.CurrentTemperature = 25.0

		repo := new(MockSectionRepository)
		repo.On("Update", 1, updateData).Return(expectedSection, nil)

		service := NewSectionDefault(repo)

		// Act
		updatedSection, err := service.Update(1, updateData)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedSection, updatedSection)
		repo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Arrange
		updateData := models.Section{ID: 999}
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = errors.New("section not found")

		repo := new(MockSectionRepository)
		repo.On("Update", 999, updateData).Return(models.Section{}, svcErr)

		service := NewSectionDefault(repo)

		// Act
		section, err := service.Update(999, updateData)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, section.ID)
		repo.AssertExpectations(t)
	})
}

func TestSectionService_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		repo := new(MockSectionRepository)
		repo.On("Delete", 1).Return(nil)

		service := NewSectionDefault(repo)

		// Act
		err := service.Delete(1)

		// Assert
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Arrange
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = errors.New("section not found")

		repo := new(MockSectionRepository)
		repo.On("Delete", 999).Return(svcErr)

		service := NewSectionDefault(repo)

		// Act
		err := service.Delete(999)

		// Assert
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
