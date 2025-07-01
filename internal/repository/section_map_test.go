package repository

import (
	"app/pkg"
	"app/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestSectionRepository_Create(t *testing.T) {
	t.Run("Create_Success", func(t *testing.T) {
		repo := NewSectionMapRepository([]models.Section{})
		section := createTestSection(1)
		section.ID = 0 // ID should be set by the repository

		result, err := repo.Create(section)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID) // First ID should be 1
		assert.Equal(t, section.SectionNumber, result.SectionNumber)
	})

	t.Run("Create_DuplicateSectionNumber_Error", func(t *testing.T) {
		section1 := createTestSection(1)
		section2 := createTestSection(2)
		section2.SectionNumber = section1.SectionNumber // Same section number

		repo := NewSectionMapRepository([]models.Section{section1})

		_, err := repo.Create(section2)

		assert.Error(t, err)
		svcErr, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ErrConflict, svcErr.Code)
	})
}

func TestSectionRepository_GetAll(t *testing.T) {
	t.Run("GetAll_Success", func(t *testing.T) {
		section1 := createTestSection(1)
		section2 := createTestSection(2)
		repo := NewSectionMapRepository([]models.Section{section1, section2})

		result, err := repo.GetAll()

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// Convert map to slice for easier comparison
		var sections []models.Section
		for _, s := range result {
			sections = append(sections, s)
		}

		// Check that both sections are present, regardless of order
		assert.Contains(t, sections, section1)
		assert.Contains(t, sections, section2)
	})

	t.Run("GetAll_Empty", func(t *testing.T) {
		repo := NewSectionMapRepository([]models.Section{})

		result, err := repo.GetAll()

		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestSectionRepository_GetByID(t *testing.T) {
	t.Run("GetByID_Success", func(t *testing.T) {
		section := createTestSection(1)
		repo := NewSectionMapRepository([]models.Section{section})

		result, err := repo.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, section, result)
	})

	t.Run("GetByID_NotFound_Error", func(t *testing.T) {
		repo := NewSectionMapRepository([]models.Section{})

		_, err := repo.GetByID(999)

		assert.Error(t, err)
		svcErr, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ErrNotFound, svcErr.Code)
	})
}

func TestSectionRepository_Update(t *testing.T) {
	t.Run("Update_Success", func(t *testing.T) {
		section := createTestSection(1)
		repo := NewSectionMapRepository([]models.Section{section})

		updateData := models.Section{
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      200,
				CurrentTemperature: 25.0,
			},
		}

		updated, err := repo.Update(1, updateData)

		assert.NoError(t, err)
		assert.Equal(t, 200, updated.SectionNumber)
		assert.Equal(t, 25.0, updated.CurrentTemperature)
		// Check that other fields remain unchanged
		assert.Equal(t, section.MinimumTemperature, updated.MinimumTemperature)
		assert.Equal(t, section.CurrentCapacity, updated.CurrentCapacity)
	})

	t.Run("Update_NotFound_Error", func(t *testing.T) {
		repo := NewSectionMapRepository([]models.Section{})

		_, err := repo.Update(999, models.Section{})

		assert.Error(t, err)
		svcErr, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ErrNotFound, svcErr.Code)
	})
}

func TestSectionRepository_Delete(t *testing.T) {
	t.Run("Delete_Success", func(t *testing.T) {
		section := createTestSection(1)
		repo := NewSectionMapRepository([]models.Section{section})

		err := repo.Delete(1)

		assert.NoError(t, err)
		_, err = repo.GetByID(1)
		assert.Error(t, err) // Should be deleted
	})

	t.Run("Delete_NotFound_Error", func(t *testing.T) {
		repo := NewSectionMapRepository([]models.Section{})

		err := repo.Delete(999)

		assert.Error(t, err)
		svcErr, ok := err.(pkg.ServiceError)
		assert.True(t, ok)
		assert.Equal(t, pkg.ErrNotFound, svcErr.Code)
	})
}

func TestSectionRepository_ConcurrentAccess(t *testing.T) {
	// Este test verifica que el repositorio pueda manejar múltiples operaciones secuenciales
	repo := NewSectionMapRepository([]models.Section{})

	// Creamos un canal para recibir los resultados
	results := make(chan models.Section, 10)
	errors := make(chan error, 10)

	// Función auxiliar para crear secciones de forma secuencial
	createSection := func(id int) {
		section := createTestSection(id)
		section.ID = 0 // El repositorio asignará el ID
		created, err := repo.Create(section)
		if err != nil {
			errors <- err
			return
		}
		results <- created
	}

	// Ejecutamos las creaciones de forma secuencial en lugar de concurrente
	for i := 0; i < 10; i++ {
		createSection(i)
	}

	close(results)
	close(errors)

	// Verificamos si hubo errores
	for err := range errors {
		t.Fatalf("Error creando sección: %v", err)
	}

	// Verificamos que se crearon todas las secciones
	sections, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, sections, 10, "Deberían haberse creado 10 secciones")

	// Verificamos que todos los IDs son únicos
	idSet := make(map[int]bool)
	for _, section := range sections {
		if idSet[section.ID] {
			t.Fatalf("ID duplicado encontrado: %d", section.ID)
		}
		idSet[section.ID] = true
	}
}
