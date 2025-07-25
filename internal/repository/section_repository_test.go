package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSectionRepository_GetAll(t *testing.T) {
	t.Run("Devuelve todas las secciones correctamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id", "section_number", "current_temperature", "current_capacity",
			"minimum_temperature", "minimum_capacity", "product_type_id", "warehouse_id",
		}).
			AddRow(1, 10, 5.0, 50, 2.0, 10, 200, 1).
			AddRow(2, 20, 7.5, 60, 3.2, 15, 201, 2)

		mock.ExpectQuery("SELECT id, section_number, current_temperature, current_capacity, minimum_temperature, minimum_capacity, product_type_id, warehouse_id FROM sections").
			WillReturnRows(rows)

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.GetAll()

		// Assert
		require.NoError(t, err)
		require.Len(t, result, 2)

		sec1 := result[1]
		require.Equal(t, 1, sec1.ID)
		require.Equal(t, 10, sec1.SectionNumber)
		require.Equal(t, 5.0, sec1.CurrentTemperature)
		require.Equal(t, 50, sec1.CurrentCapacity)
		require.Equal(t, 2.0, sec1.MinimumTemperature)
		require.Equal(t, 10, sec1.MinimumCapacity)
		require.Equal(t, 200, sec1.ProductTypeID)
		require.Equal(t, 1, sec1.WarehouseID)

		sec2 := result[2]
		require.Equal(t, 2, sec2.ID)
		require.Equal(t, 20, sec2.SectionNumber)
		require.Equal(t, 7.5, sec2.CurrentTemperature)
		require.Equal(t, 60, sec2.CurrentCapacity)
		require.Equal(t, 3.2, sec2.MinimumTemperature)
		require.Equal(t, 15, sec2.MinimumCapacity)
		require.Equal(t, 201, sec2.ProductTypeID)
		require.Equal(t, 2, sec2.WarehouseID)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Devuelve error interno si falla la query", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, section_number, current_temperature, current_capacity, minimum_temperature, minimum_capacity, product_type_id, warehouse_id FROM sections").
			WillReturnError(assert.AnError)

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.GetAll()

		// Assert
		require.Error(t, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Devuelve error not found si no hay filas", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, section_number, current_temperature, current_capacity, minimum_temperature, minimum_capacity, product_type_id, warehouse_id FROM sections").
			WillReturnError(sql.ErrNoRows)

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.GetAll()

		// Assert
		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("no se encontraron sections")
		require.Equal(t, svcErr, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSectionRepository_GetByID(t *testing.T) {
	t.Run("Devuelve la seccion correctamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, section_number, current_temperature, current_capacity, minimum_temperature, minimum_capacity, product_type_id, warehouse_id FROM sections WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "section_number", "current_temperature", "current_capacity", "minimum_temperature", "minimum_capacity", "product_type_id", "warehouse_id"}).
				AddRow(1, 10, 5.0, 50, 2.0, 10, 200, 1))

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.GetByID(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, 10, result.SectionNumber)
		require.Equal(t, 5.0, result.CurrentTemperature)
		require.Equal(t, 50, result.CurrentCapacity)
		require.Equal(t, 2.0, result.MinimumTemperature)
		require.Equal(t, 10, result.MinimumCapacity)
		require.Equal(t, 200, result.ProductTypeID)
		require.Equal(t, 1, result.WarehouseID)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Devuelve error interno si falla la query", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, section_number, current_temperature, current_capacity, minimum_temperature, minimum_capacity, product_type_id, warehouse_id FROM sections WHERE id = ?").
			WithArgs(1).
			WillReturnError(assert.AnError)

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.GetByID(1)

		// Assert
		require.Error(t, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Devuelve error not found si no se encuentra la seccion", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, section_number, current_temperature, current_capacity, minimum_temperature, minimum_capacity, product_type_id, warehouse_id FROM sections WHERE id = ?").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.GetByID(1)

		// Assert
		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("section con id 1 no encontrada")
		require.Equal(t, svcErr, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSectionRepository_Create(t *testing.T) {
	t.Run("Se crea correctamente la seccion", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("INSERT INTO sections.*VALUES.*").
			WithArgs(10, 5.0, 50, 2.0, 10, 200, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.Create(models.Section{
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      10,
				CurrentTemperature: 5.0,
				CurrentCapacity:    50,
				MinimumTemperature: 2.0,
				MinimumCapacity:    10,
				ProductTypeID:      200,
				WarehouseID:        1,
			},
		})

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Devuelve error interno si falla la query", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("INSERT INTO sections.*VALUES.*").
			WithArgs(10, 5.0, 50, 2.0, 10, 200, 1).
			WillReturnError(assert.AnError)

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.Create(models.Section{
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      10,
				CurrentTemperature: 5.0,
				CurrentCapacity:    50,
				MinimumTemperature: 2.0,
				MinimumCapacity:    10,
				ProductTypeID:      200,
				WarehouseID:        1,
			},
		})

		// Assert
		require.Error(t, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSectionRepository_Update(t *testing.T) {
	t.Run("Se actualiza correctamente la seccion", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("UPDATE sections SET section_number = ?, current_temperature = ?, current_capacity = ?, minimum_temperature = ?, minimum_capacity = ?, product_type_id = ?, warehouse_id = ? WHERE id = ?").
			WithArgs(10, 5.0, 50, 2.0, 10, 200, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.Update(1, models.Section{
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      10,
				CurrentTemperature: 5.0,
				CurrentCapacity:    50,
				MinimumTemperature: 2.0,
				MinimumCapacity:    10,
				ProductTypeID:      200,
				WarehouseID:        1,
			},
		})

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Devuelve error interno si falla la query", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("UPDATE sections SET section_number = ?, current_temperature = ?, current_capacity = ?, minimum_temperature = ?, minimum_capacity = ?, product_type_id = ?, warehouse_id = ? WHERE id = ?").
			WithArgs(10, 5.0, 50, 2.0, 10, 200, 1, 1).
			WillReturnError(assert.AnError)

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.Update(1, models.Section{
			ID: 1,
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      10,
				CurrentTemperature: 5.0,
				CurrentCapacity:    50,
				MinimumTemperature: 2.0,
				MinimumCapacity:    10,
				ProductTypeID:      200,
				WarehouseID:        1,
			},
		})

		// Assert
		require.Error(t, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Devuelve error not found si no se encuentra la seccion", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("UPDATE sections SET section_number = ?, current_temperature = ?, current_capacity = ?, minimum_temperature = ?, minimum_capacity = ?, product_type_id = ?, warehouse_id = ? WHERE id = ?").
			WithArgs(10, 5.0, 50, 2.0, 10, 200, 1, 1).
			WillReturnError(sql.ErrNoRows)

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.Update(1, models.Section{
			ID: 1,
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      10,
				CurrentTemperature: 5.0,
				CurrentCapacity:    50,
				MinimumTemperature: 2.0,
				MinimumCapacity:    10,
				ProductTypeID:      200,
				WarehouseID:        1,
			},
		})

		// Assert
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("section con id 1 no encontrada")
		require.Equal(t, svcErr, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSectionRepository_Delete(t *testing.T) {
	t.Run("Se elimina correctamente la seccion", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("DELETE FROM sections WHERE id = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewSectionSqlRepository(db)

		// Act
		err = repo.Delete(1)

		// Assert
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Devuelve error interno si falla la query", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("DELETE FROM sections WHERE id = ?").
			WithArgs(1).
			WillReturnError(assert.AnError)

		repo := NewSectionSqlRepository(db)

		// Act
		err = repo.Delete(1)

		// Assert
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Devuelve error not found si no se encuentra la seccion", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("DELETE FROM sections WHERE id = ?").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		repo := NewSectionSqlRepository(db)

		// Act
		err = repo.Delete(1)

		// Assert
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
