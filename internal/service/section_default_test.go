package service

import (
	"testing"

	"app/pkg"
	"app/pkg/models"
	sectionmock "app/test/section"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("repository_error", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		repo.On("GetAll").Return(map[int]models.Section{}, pkg.ServiceErrors[pkg.ErrInternalServer])
		service := NewSectionDefault(repo)

		// Act
		result, err := service.GetAll()

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		assert.Empty(t, result)
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

	t.Run("not_found", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		sectionID := 999
		repo.On("GetByID", sectionID).Return(models.Section{}, pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewSectionDefault(repo)

		// Act
		result, err := service.GetByID(sectionID)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		assert.Equal(t, models.Section{}, result)
		repo.AssertExpectations(t)
	})

	t.Run("internal_error", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		sectionID := 1
		repo.On("GetByID", sectionID).Return(models.Section{}, pkg.ServiceErrors[pkg.ErrInternalServer])
		service := NewSectionDefault(repo)

		// Act
		result, err := service.GetByID(sectionID)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		assert.Equal(t, models.Section{}, result)
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

	t.Run("conflict_error", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		input := models.Section{ID: 3}
		repo.On("Create", input).Return(models.Section{}, pkg.ServiceErrors[pkg.ErrConflict])
		service := NewSectionDefault(repo)

		// Act
		result, err := service.Create(input)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict], err)
		assert.Equal(t, models.Section{}, result)
		repo.AssertExpectations(t)
	})

	t.Run("internal_error", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		input := models.Section{ID: 3}
		repo.On("Create", input).Return(models.Section{}, pkg.ServiceErrors[pkg.ErrInternalServer])
		service := NewSectionDefault(repo)

		// Act
		result, err := service.Create(input)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		assert.Equal(t, models.Section{}, result)
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

	t.Run("get_by_id_error", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		id := 999
		upd := models.Section{SectionAttributes: models.SectionAttributes{SectionNumber: 2}}
		repo.On("GetByID", id).Return(models.Section{}, pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewSectionDefault(repo)

		// Act
		result, err := service.Update(id, upd)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		assert.Equal(t, models.Section{}, result)
		repo.AssertExpectations(t)
	})

	t.Run("update_error", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		id := 4
		old := models.Section{ID: id, SectionAttributes: models.SectionAttributes{SectionNumber: 1}}
		upd := models.Section{SectionAttributes: models.SectionAttributes{SectionNumber: 2}}
		updated := models.Section{ID: id, SectionAttributes: models.SectionAttributes{SectionNumber: 2}}
		repo.On("GetByID", id).Return(old, nil)
		repo.On("Update", id, updated).Return(models.Section{}, pkg.ServiceErrors[pkg.ErrInternalServer])
		service := NewSectionDefault(repo)

		// Act
		result, err := service.Update(id, upd)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		assert.Equal(t, models.Section{}, result)
		repo.AssertExpectations(t)
	})

	t.Run("update_all_fields", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		id := 5
		old := models.Section{
			ID: id,
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      1,
				CurrentTemperature: 10,
				MinimumTemperature: 5,
				CurrentCapacity:    100,
				MinimumCapacity:    50,
				MaximumCapacity:    200,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
		}
		upd := models.Section{
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      2,
				CurrentTemperature: 15,
				MinimumTemperature: 8,
				CurrentCapacity:    120,
				MinimumCapacity:    60,
				MaximumCapacity:    250,
				WarehouseID:        2,
				ProductTypeID:      2,
			},
		}
		updated := models.Section{
			ID: id,
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      2,
				CurrentTemperature: 15,
				MinimumTemperature: 8,
				CurrentCapacity:    120,
				MinimumCapacity:    60,
				MaximumCapacity:    250,
				WarehouseID:        2,
				ProductTypeID:      2,
			},
		}
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

	t.Run("delete_non_existent", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		sectionID := 999
		repo.On("Delete", sectionID).Return(pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewSectionDefault(repo)

		// Act
		err := service.Delete(sectionID)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		repo.AssertExpectations(t)
	})

	t.Run("delete_internal_error", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		sectionID := 10
		repo.On("Delete", sectionID).Return(pkg.ServiceErrors[pkg.ErrInternalServer])
		service := NewSectionDefault(repo)

		// Act
		err := service.Delete(sectionID)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
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

	t.Run("all_sections_error", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		repo.On("GetReportProductsAllSections").Return([]models.SectionReport{}, pkg.ServiceErrors[pkg.ErrInternalServer])
		service := NewSectionDefault(repo)

		// Act
		result, err := service.ReportProductsBySection(nil)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		assert.Empty(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("specific_section_error", func(t *testing.T) {
		// Arrange
		repo := new(sectionmock.MockSectionRepository)
		id := 999
		repo.On("GetReportProductsBySection", id).Return([]models.SectionReport{}, pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewSectionDefault(repo)

		// Act
		result, err := service.ReportProductsBySection(&id)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		assert.Empty(t, result)
		repo.AssertExpectations(t)
	})
}
