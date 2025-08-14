package repository

import (
	"app/pkg"
	"app/pkg/models"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestCarrySql_Create(t *testing.T) {
	t.Run("success - create carry", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		carry := models.Carry{
			Cid:         "CAR001",
			CompanyName: "Transport Co",
			Address:     "123 Main St",
			Telephone:   "555-0101",
			LocalityId:  1,
		}

		mock.ExpectExec("INSERT INTO carries \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("CAR001", "Transport Co", "123 Main St", "555-0101", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		result, err := repo.Create(carry)

		// Assert
		require.NoError(t, err)
		require.Equal(t, carry, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - duplicate cid", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		carry := models.Carry{
			Cid:         "CAR001",
			CompanyName: "Transport Co",
			Address:     "123 Main St",
			Telephone:   "555-0101",
			LocalityId:  1,
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1062,
			Message: "Duplicate entry 'CAR001' for key 'cid'",
		}

		mock.ExpectExec("INSERT INTO carries \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("CAR001", "Transport Co", "123 Main St", "555-0101", 1).
			WillReturnError(mysqlErr)

		// Act
		result, err := repo.Create(carry)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Carry{}, result)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ErrConflict, serviceErr.Code)
		require.Equal(t, "cid already exists", serviceErr.InternalError.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - locality not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		carry := models.Carry{
			Cid:         "CAR001",
			CompanyName: "Transport Co",
			Address:     "123 Main St",
			Telephone:   "555-0101",
			LocalityId:  999,
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails",
		}

		mock.ExpectExec("INSERT INTO carries \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("CAR001", "Transport Co", "123 Main St", "555-0101", 999).
			WillReturnError(mysqlErr)

		// Act
		result, err := repo.Create(carry)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Carry{}, result)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ErrConflict, serviceErr.Code)
		require.Equal(t, "locality not found", serviceErr.InternalError.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - other mysql error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		carry := models.Carry{
			Cid:         "CAR001",
			CompanyName: "Transport Co",
			Address:     "123 Main St",
			Telephone:   "555-0101",
			LocalityId:  1,
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1064,
			Message: "You have an error in your SQL syntax",
		}

		mock.ExpectExec("INSERT INTO carries \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("CAR001", "Transport Co", "123 Main St", "555-0101", 1).
			WillReturnError(mysqlErr)

		// Act
		result, err := repo.Create(carry)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Carry{}, result)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ErrInternalServer, serviceErr.Code)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - non-mysql error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		carry := models.Carry{
			Cid:         "CAR001",
			CompanyName: "Transport Co",
			Address:     "123 Main St",
			Telephone:   "555-0101",
			LocalityId:  1,
		}

		mock.ExpectExec("INSERT INTO carries \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("CAR001", "Transport Co", "123 Main St", "555-0101", 1).
			WillReturnError(errors.New("generic database error"))

		// Act
		result, err := repo.Create(carry)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Carry{}, result)
		require.Equal(t, "generic database error", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCarrySql_SearchByLocality(t *testing.T) {
	t.Run("success - search carries by specific locality", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		expectedResult := map[string]models.CarryByLocality{
			"1": {
				LocalityId:   "1",
				LocalityName: "Buenos Aires",
				CarriesCount: 3,
			},
		}

		rows := sqlmock.NewRows([]string{"locality_id", "locality_name", "carries_count"}).
			AddRow("1", "Buenos Aires", 3)

		mock.ExpectQuery("SELECT l.id AS locality_id, l.locality_name, COUNT\\(c.id\\) AS carries_count FROM localities l LEFT JOIN carries c ON l.id = c.locality_id where l.id = \\? GROUP BY l.id, l.locality_name").
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		result, err := repo.SearchByLocality(1)

		// Assert
		require.NoError(t, err)
		require.Len(t, result, 1)
		require.Equal(t, expectedResult["1"], result["1"])
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - search all carries by locality", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		expectedResult := map[string]models.CarryByLocality{
			"1": {
				LocalityId:   "1",
				LocalityName: "Buenos Aires",
				CarriesCount: 3,
			},
			"2": {
				LocalityId:   "2",
				LocalityName: "Córdoba",
				CarriesCount: 2,
			},
		}

		rows := sqlmock.NewRows([]string{"locality_id", "locality_name", "carries_count"}).
			AddRow("1", "Buenos Aires", 3).
			AddRow("2", "Córdoba", 2)

		mock.ExpectQuery("SELECT l.id AS locality_id, l.locality_name, COUNT\\(c.id\\) AS carries_count FROM localities l LEFT JOIN carries c ON l.id = c.locality_id  GROUP BY l.id, l.locality_name").
			WillReturnRows(rows)

		// Act
		result, err := repo.SearchByLocality(-1)

		// Assert
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, expectedResult["1"], result["1"])
		require.Equal(t, expectedResult["2"], result["2"])
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - empty result for specific locality", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		rows := sqlmock.NewRows([]string{"locality_id", "locality_name", "carries_count"})

		mock.ExpectQuery("SELECT l.id AS locality_id, l.locality_name, COUNT\\(c.id\\) AS carries_count FROM localities l LEFT JOIN carries c ON l.id = c.locality_id where l.id = \\? GROUP BY l.id, l.locality_name").
			WithArgs(999).
			WillReturnRows(rows)

		// Act
		result, err := repo.SearchByLocality(999)

		// Assert
		require.NoError(t, err)
		require.Len(t, result, 0)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails for specific locality", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		mock.ExpectQuery("SELECT l.id AS locality_id, l.locality_name, COUNT\\(c.id\\) AS carries_count FROM localities l LEFT JOIN carries c ON l.id = c.locality_id where l.id = \\? GROUP BY l.id, l.locality_name").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		// Act
		result, err := repo.SearchByLocality(1)

		// Assert
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, "database error", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails for all localities", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		mock.ExpectQuery("SELECT l.id AS locality_id, l.locality_name, COUNT\\(c.id\\) AS carries_count FROM localities l LEFT JOIN carries c ON l.id = c.locality_id  GROUP BY l.id, l.locality_name").
			WillReturnError(errors.New("database error"))

		// Act
		result, err := repo.SearchByLocality(-1)

		// Assert
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, "database error", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - scan fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCarrySql(db)

		rows := sqlmock.NewRows([]string{"locality_id", "locality_name", "carries_count"}).
			AddRow("invalid", "Buenos Aires", "invalid_count")

		mock.ExpectQuery("SELECT l.id AS locality_id, l.locality_name, COUNT\\(c.id\\) AS carries_count FROM localities l LEFT JOIN carries c ON l.id = c.locality_id where l.id = \\? GROUP BY l.id, l.locality_name").
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		result, err := repo.SearchByLocality(1)

		// Assert
		require.Error(t, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNewCarrySql(t *testing.T) {
	t.Run("success - create new carry sql repository", func(t *testing.T) {
		// Arrange
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Act
		repo := NewCarrySql(db)

		// Assert
		require.NotNil(t, repo)
		require.Equal(t, db, repo.db)
	})
}
