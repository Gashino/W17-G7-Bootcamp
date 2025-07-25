package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
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
			AddRow(1, 10, 5, 50, 2, 10, 200, 1).
			AddRow(2, 20, 7, 60, 3, 15, 201, 2)

		mock.ExpectQuery("SELECT id, section_number, current_temperature, current_capacity, minimum_temperature, minimum_capacity, product_type_id, warehouse_id FROM sections").
			WillReturnRows(rows)

		repo := NewSectionSqlRepository(db)

		// Act
		result, err := repo.GetAll()

		// Assert
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, 10, result[1].SectionNumber)
		require.Equal(t, 20, result[2].SectionNumber)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSectionRepository_GetByID(t *testing.T) {

}

func TestSectionRepository_Create(t *testing.T) {

}

func TestSectionRepository_Update(t *testing.T) {

}

func TestSectionRepository_Delete(t *testing.T) {

}
