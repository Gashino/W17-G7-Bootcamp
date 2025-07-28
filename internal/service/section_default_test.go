package service

import (
	"testing"

	"app/pkg/models"
	sectionmock "app/test/section"

	"github.com/stretchr/testify/assert"
)

func TestSectionDefault_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		expected := map[int]models.Section{1: {ID: 1}}
		repo.On("GetAll").Return(expected, nil)
		service := NewSectionDefault(repo)

		// Act
		result, err := service.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})
}

func TestSectionDefault_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		expected := models.Section{ID: 2}
		repo.On("GetByID", 2).Return(expected, nil)
		service := NewSectionDefault(repo)

		// Act
		result, err := service.GetByID(2)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})
}

func TestSectionDefault_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		input := models.Section{ID: 3}
		repo.On("Create", input).Return(input, nil)
		service := NewSectionDefault(repo)

		// Act
		result, err := service.Create(input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, input, result)
		repo.AssertExpectations(t)
	})
}

func TestSectionDefault_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		id := 4
		old := models.Section{ID: id, SectionAttributes: models.SectionAttributes{SectionNumber: 1}}
		upd := models.Section{SectionAttributes: models.SectionAttributes{SectionNumber: 2}}
		updated := models.Section{ID: id, SectionAttributes: models.SectionAttributes{SectionNumber: 2}}
		repo.On("GetByID", id).Return(old, nil)
		repo.On("Update", id, updated).Return(updated, nil)
		service := NewSectionDefault(repo)

		// Act
		result, err := service.Update(id, upd)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, updated, result)
		repo.AssertExpectations(t)
	})
}

func TestSectionDefault_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		repo.On("Delete", 5).Return(nil)
		service := NewSectionDefault(repo)

		// Act
		err := service.Delete(5)

		// Assert
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestSectionDefault_ReportProductsBySection(t *testing.T) {
	t.Run("all sections", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		reports := []models.SectionReport{{SectionId: 1}}
		repo.On("GetReportProductsAllSections").Return(reports, nil)
		service := NewSectionDefault(repo)

		// Act
		result, err := service.ReportProductsBySection(nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, reports, result)
		repo.AssertExpectations(t)
	})

	t.Run("specific section", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		reports := []models.SectionReport{{SectionId: 1}}
		id := 1
		repo.On("GetReportProductsBySection", id).Return(reports, nil)
		service := NewSectionDefault(repo)

		// Act
		result, err := service.ReportProductsBySection(&id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, reports, result)
		repo.AssertExpectations(t)
	})
}
