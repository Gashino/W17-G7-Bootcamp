package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestLocalityRepository_GetById(t *testing.T) {
	t.Run("Cuando encuentra la locality retorna los datos correctamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id", "locality_name", "province_name", "country_name",
		}).AddRow(1, "Buenos Aires", "Buenos Aires", "Argentina")

		mock.ExpectQuery("SELECT id, locality_name, province_name, country_name FROM localities WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(rows)

		repo := NewLocalitySql(db)

		// Act
		result, err := repo.GetById(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "Buenos Aires", result.LocalityName)
		require.Equal(t, "Buenos Aires", result.ProvinceName)
		require.Equal(t, "Argentina", result.CountryName)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando no encuentra la locality retorna error not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, locality_name, province_name, country_name FROM localities WHERE id = \\?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		repo := NewLocalitySql(db)

		// Act
		result, err := repo.GetById(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Locality{}, result)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error de base de datos retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, locality_name, province_name, country_name FROM localities WHERE id = \\?").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		repo := NewLocalitySql(db)

		// Act
		result, err := repo.GetById(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Locality{}, result)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLocalityRepository_Create(t *testing.T) {
	t.Run("Cuando los datos son válidos crea la locality exitosamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		inputLocality := models.Locality{
			LocalitiesAttributes: models.LocalitiesAttributes{
				LocalityName: "Córdoba",
				ProvinceName: "Córdoba",
				CountryName:  "Argentina",
			},
		}

		// Mock para verificar que no existe locality con el mismo nombre
		mock.ExpectQuery("SELECT id FROM localities WHERE locality_name = \\?").
			WithArgs("Córdoba").
			WillReturnError(sql.ErrNoRows)

		// Mock para la inserción
		mock.ExpectExec("INSERT INTO localities \\(locality_name, province_name, country_name\\) VALUES \\(\\?, \\?, \\?\\)").
			WithArgs("Córdoba", "Córdoba", "Argentina").
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewLocalitySql(db)

		// Act
		result, err := repo.Create(inputLocality)

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "Córdoba", result.LocalityName)
		require.Equal(t, "Córdoba", result.ProvinceName)
		require.Equal(t, "Argentina", result.CountryName)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando la locality ya existe retorna error de conflicto", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		inputLocality := models.Locality{
			LocalitiesAttributes: models.LocalitiesAttributes{
				LocalityName: "Buenos Aires",
				ProvinceName: "Buenos Aires",
				CountryName:  "Argentina",
			},
		}

		// Mock para verificar que existe locality con el mismo nombre
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("SELECT id FROM localities WHERE locality_name = \\?").
			WithArgs("Buenos Aires").
			WillReturnRows(rows)

		repo := NewLocalitySql(db)

		// Act
		result, err := repo.Create(inputLocality)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Locality{}, result)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error en la inserción retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		inputLocality := models.Locality{
			LocalitiesAttributes: models.LocalitiesAttributes{
				LocalityName: "Mendoza",
				ProvinceName: "Mendoza",
				CountryName:  "Argentina",
			},
		}

		// Mock para verificar que no existe locality con el mismo nombre
		mock.ExpectQuery("SELECT id FROM localities WHERE locality_name = \\?").
			WithArgs("Mendoza").
			WillReturnError(sql.ErrNoRows)

		// Mock para la inserción con error
		mock.ExpectExec("INSERT INTO localities \\(locality_name, province_name, country_name\\) VALUES \\(\\?, \\?, \\?\\)").
			WithArgs("Mendoza", "Mendoza", "Argentina").
			WillReturnError(sql.ErrConnDone)

		repo := NewLocalitySql(db)

		// Act
		result, err := repo.Create(inputLocality)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Locality{}, result)
		require.Contains(t, err.Error(), "failed to create Locality")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLocalityRepository_GetCantSellersByLocality(t *testing.T) {
	t.Run("Cuando encuentra la locality con sellers retorna el reporte correctamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"idLocalidad", "nombreLocalidad", "cantidadSellers",
		}).AddRow(1, "Buenos Aires", "5")

		mock.ExpectQuery("SELECT l.id AS idLocalidad, l.locality_name AS nombreLocalidad, COUNT\\(s.id\\) AS cantidadSellers").
			WithArgs(1).
			WillReturnRows(rows)

		repo := NewLocalitySql(db)

		// Act
		result, err := repo.GetCantSellersByLocality(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.NotNil(t, result.LocalityName)
		require.Equal(t, "Buenos Aires", *result.LocalityName)
		require.NotNil(t, result.SellerCount)
		require.Equal(t, "5", *result.SellerCount)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando la locality no existe retorna error not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT l.id AS idLocalidad, l.locality_name AS nombreLocalidad, COUNT\\(s.id\\) AS cantidadSellers").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		repo := NewLocalitySql(db)

		// Act
		result, err := repo.GetCantSellersByLocality(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.LocalityBySellerResponse{}, result)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error de base de datos retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT l.id AS idLocalidad, l.locality_name AS nombreLocalidad, COUNT\\(s.id\\) AS cantidadSellers").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		repo := NewLocalitySql(db)

		// Act
		result, err := repo.GetCantSellersByLocality(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.LocalityBySellerResponse{}, result)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando la locality existe pero no tiene sellers retorna cero", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"idLocalidad", "nombreLocalidad", "cantidadSellers",
		}).AddRow(1, "Tucumán", "0")

		mock.ExpectQuery("SELECT l.id AS idLocalidad, l.locality_name AS nombreLocalidad, COUNT\\(s.id\\) AS cantidadSellers").
			WithArgs(1).
			WillReturnRows(rows)

		repo := NewLocalitySql(db)

		// Act
		result, err := repo.GetCantSellersByLocality(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.NotNil(t, result.LocalityName)
		require.Equal(t, "Tucumán", *result.LocalityName)
		require.NotNil(t, result.SellerCount)
		require.Equal(t, "0", *result.SellerCount)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
