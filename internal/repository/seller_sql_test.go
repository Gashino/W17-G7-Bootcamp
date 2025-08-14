package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestSellerRepository_FindAll(t *testing.T) {
	t.Run("Devuelve todos los sellers correctamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).
			AddRow(1, "12345", "Company One", "Address One", "123456789", 1).
			AddRow(2, "67890", "Company Two", "Address Two", "987654321", 2)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers").
			WillReturnRows(rows)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.FindAll()

		// Assert
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, "12345", result[1].CId)
		require.Equal(t, "Company One", result[1].CompanyName)
		require.Equal(t, "67890", result[2].CId)
		require.Equal(t, "Company Two", result[2].CompanyName)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error de base de datos retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers").
			WillReturnError(sql.ErrConnDone)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.FindAll()

		// Assert
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error en el scan de filas retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Crear filas con datos incompatibles para causar error en Scan
		rows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow("invalid_id", "12345", "Company One", "Address One", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers").
			WillReturnRows(rows)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.FindAll()

		// Assert
		require.Error(t, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepository_GetById(t *testing.T) {
	t.Run("Cuando encuentra el seller retorna los datos correctamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(1, "12345", "Test Company", "Test Address", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(rows)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.GetById(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "12345", result.CId)
		require.Equal(t, "Test Company", result.CompanyName)
		require.Equal(t, "Test Address", result.Address)
		require.Equal(t, "123456789", result.Telephone)
		require.Equal(t, 1, result.LocalityID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando no encuentra el seller retorna error not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.GetById(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error de base de datos retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.GetById(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepository_Create(t *testing.T) {
	t.Run("Cuando los datos son válidos crea el seller exitosamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		inputSeller := models.Seller{
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "New Company",
				Address:     "New Address",
				Telephone:   "123456789",
				LocalityID:  1,
			},
		}

		// Mock para verificar que no existe seller con el mismo CId
		mock.ExpectQuery("SELECT id FROM sellers WHERE cid = \\?").
			WithArgs("12345").
			WillReturnError(sql.ErrNoRows)

		// Mock para la inserción
		mock.ExpectExec("INSERT INTO sellers \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("12345", "New Company", "New Address", "123456789", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewSellerSql(db)

		// Act
		result, err := repo.Create(inputSeller)

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "12345", result.CId)
		require.Equal(t, "New Company", result.CompanyName)
		require.Equal(t, "New Address", result.Address)
		require.Equal(t, "123456789", result.Telephone)
		require.Equal(t, 1, result.LocalityID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando el CId ya existe retorna error de conflicto", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		inputSeller := models.Seller{
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "Duplicate Company",
				Address:     "Some Address",
				Telephone:   "123456789",
				LocalityID:  1,
			},
		}

		// Mock para verificar que existe seller con el mismo CId
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("SELECT id FROM sellers WHERE cid = \\?").
			WithArgs("12345").
			WillReturnRows(rows)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.Create(inputSeller)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error en la inserción retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		inputSeller := models.Seller{
			SellerAttributes: models.SellerAttributes{
				CId:         "54321",
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "987654321",
				LocalityID:  1,
			},
		}

		// Mock para verificar que no existe seller con el mismo CId
		mock.ExpectQuery("SELECT id FROM sellers WHERE cid = \\?").
			WithArgs("54321").
			WillReturnError(sql.ErrNoRows)

		// Mock para la inserción con error
		mock.ExpectExec("INSERT INTO sellers \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("54321", "Test Company", "Test Address", "987654321", 1).
			WillReturnError(sql.ErrConnDone)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.Create(inputSeller)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Contains(t, err.Error(), "failed to create seller")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error de base de datos en verificación de CId retorna error not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		inputSeller := models.Seller{
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "123456789",
				LocalityID:  1,
			},
		}

		// Mock para verificar CId con error de conexión (no sql.ErrNoRows)
		mock.ExpectQuery("SELECT id FROM sellers WHERE cid = \\?").
			WithArgs("12345").
			WillReturnError(sql.ErrConnDone)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.Create(inputSeller)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error al obtener LastInsertId retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		inputSeller := models.Seller{
			SellerAttributes: models.SellerAttributes{
				CId:         "99999",
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "123456789",
				LocalityID:  1,
			},
		}

		// Mock para verificar que no existe seller con el mismo CId
		mock.ExpectQuery("SELECT id FROM sellers WHERE cid = \\?").
			WithArgs("99999").
			WillReturnError(sql.ErrNoRows)

		// Mock para la inserción exitosa pero con error en LastInsertId
		mock.ExpectExec("INSERT INTO sellers \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("99999", "Test Company", "Test Address", "123456789", 1).
			WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

		repo := NewSellerSql(db)

		// Act
		result, err := repo.Create(inputSeller)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Contains(t, err.Error(), "failed to get last insert id")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepository_UpdateFields(t *testing.T) {
	t.Run("Cuando el seller existe actualiza los campos exitosamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		updateData := models.SellerCreateRequest{
			CompanyName: models.StringPtr("Updated Company"),
			Address:     models.StringPtr("Updated Address"),
		}

		// Mock para GetById inicial (verificar que existe)
		existingRows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(1, "12345", "Old Company", "Old Address", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(existingRows)

		// Mock para la actualización
		mock.ExpectExec("UPDATE sellers SET").
			WithArgs(nil, "Updated Company", "Updated Address", nil, nil, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Mock para GetById final (retornar seller actualizado)
		updatedRows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(1, "12345", "Updated Company", "Updated Address", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(updatedRows)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.UpdateFields(1, updateData)

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "12345", result.CId)
		require.Equal(t, "Updated Company", result.CompanyName)
		require.Equal(t, "Updated Address", result.Address)
		require.Equal(t, "123456789", result.Telephone)
		require.Equal(t, 1, result.LocalityID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando el seller no existe retorna error not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		updateData := models.SellerCreateRequest{
			CompanyName: models.StringPtr("Updated Company"),
		}

		// Mock para GetById (seller no existe)
		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.UpdateFields(999, updateData)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error en la actualización retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		updateData := models.SellerCreateRequest{
			CompanyName: models.StringPtr("Updated Company"),
		}

		// Mock para GetById inicial (verificar que existe)
		existingRows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(1, "12345", "Old Company", "Old Address", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(existingRows)

		// Mock para la actualización con error
		mock.ExpectExec("UPDATE sellers SET").
			WithArgs(nil, "Updated Company", nil, nil, nil, 1).
			WillReturnError(sql.ErrConnDone)

		repo := NewSellerSql(db)

		// Act
		result, err := repo.UpdateFields(1, updateData)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Contains(t, err.Error(), "failed to update seller")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error al obtener RowsAffected retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		updateData := models.SellerCreateRequest{
			CompanyName: models.StringPtr("Updated Company"),
		}

		// Mock para GetById inicial (verificar que existe)
		existingRows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(1, "12345", "Old Company", "Old Address", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(existingRows)

		// Mock para la actualización con error en RowsAffected
		mock.ExpectExec("UPDATE sellers SET").
			WithArgs(nil, "Updated Company", nil, nil, nil, 1).
			WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

		repo := NewSellerSql(db)

		// Act
		result, err := repo.UpdateFields(1, updateData)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Contains(t, err.Error(), "failed to get rows affected")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando no se afectan filas retorna error not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		updateData := models.SellerCreateRequest{
			CompanyName: models.StringPtr("Updated Company"),
		}

		// Mock para GetById inicial (verificar que existe)
		existingRows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(1, "12345", "Old Company", "Old Address", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(existingRows)

		// Mock para la actualización sin afectar filas
		mock.ExpectExec("UPDATE sellers SET").
			WithArgs(nil, "Updated Company", nil, nil, nil, 1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := NewSellerSql(db)

		// Act
		result, err := repo.UpdateFields(1, updateData)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepository_DeleteSeller(t *testing.T) {
	t.Run("Cuando el seller existe lo elimina exitosamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Mock para GetById (verificar que existe)
		existingRows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(1, "12345", "Test Company", "Test Address", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(existingRows)

		// Mock para la eliminación
		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repo := NewSellerSql(db)

		// Act
		err = repo.DeleteSeller(1)

		// Assert
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando el seller no existe retorna error not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Mock para GetById (seller no existe)
		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		repo := NewSellerSql(db)

		// Act
		err = repo.DeleteSeller(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error en la eliminación retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Mock para GetById (verificar que existe)
		existingRows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(1, "12345", "Test Company", "Test Address", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(existingRows)

		// Mock para la eliminación con error
		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		repo := NewSellerSql(db)

		// Act
		err = repo.DeleteSeller(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando hay error al obtener RowsAffected retorna el error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Mock para GetById (verificar que existe)
		existingRows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(1, "12345", "Test Company", "Test Address", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(existingRows)

		// Mock para la eliminación con error en RowsAffected
		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

		repo := NewSellerSql(db)

		// Act
		err = repo.DeleteSeller(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Cuando no se afectan filas retorna error not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Mock para GetById (verificar que existe)
		existingRows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(1, "12345", "Test Company", "Test Address", "123456789", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(existingRows)

		// Mock para la eliminación sin afectar filas
		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := NewSellerSql(db)

		// Act
		err = repo.DeleteSeller(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
