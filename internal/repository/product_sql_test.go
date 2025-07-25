package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
)

// Helper functions para crear punteros
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestProductSql_GetAll(t *testing.T) {
	t.Run("success_with_multiple_products", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "net_weight", "expiration_rate",
			"recommended_freezing_temperature", "freezing_rate", "product_type_id",
			"seller_id", "width", "height", "length",
		}).
			AddRow(1, "PROD001", "Product 1", 10.5, 30, -5.5, 10, 1, 1, 20.5, 10.5, 30.5).
			AddRow(2, "PROD002", "Product 2", 5.5, 15, -10.5, 5, 2, 2, 15.5, 5.5, 25.5)

		mock.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)

		repo := ProductSql{db: db}

		// Act
		products := repo.GetAll()

		// Assert
		require.NotNil(t, products)
		require.Equal(t, 2, len(products))

		// Check first product
		product1, exists := products[1]
		require.True(t, exists)
		require.Equal(t, 1, product1.ID)
		require.Equal(t, "PROD001", *product1.ProductCode)
		require.Equal(t, "Product 1", *product1.Description)
		require.Equal(t, 10.5, *product1.NetWeight)
		require.Equal(t, 30, *product1.ExpirationRate)
		require.Equal(t, -5.5, *product1.RecommendedFreezingTemperature)
		require.Equal(t, 10, *product1.FreezingRate)
		require.Equal(t, 1, *product1.ProductTypeId)
		require.Equal(t, 1, *product1.SellerId)
		require.Equal(t, 20.5, *product1.Width)
		require.Equal(t, 10.5, *product1.Height)
		require.Equal(t, 30.5, *product1.Length)

		// Check second product
		product2, exists := products[2]
		require.True(t, exists)
		require.Equal(t, 2, product2.ID)
		require.Equal(t, "PROD002", *product2.ProductCode)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM products").WillReturnError(sql.ErrConnDone)

		repo := ProductSql{db: db}

		// Act
		products := repo.GetAll()

		// Assert
		require.Nil(t, products)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Create rows with mismatched column types to cause scan error
		rows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "net_weight", "expiration_rate",
			"recommended_freezing_temperature", "freezing_rate", "product_type_id",
			"seller_id", "width", "height", "length",
		}).
			// Add a row with a string where an int is expected
			AddRow("not_an_int", "PROD001", "Product 1", 10.5, 30, -5.5, 10, 1, 1, 20.5, 10.5, 30.5)

		mock.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)

		repo := ProductSql{db: db}

		// Act
		products := repo.GetAll()

		// Assert
		require.NotNil(t, products)
		require.Equal(t, 0, len(products))

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductSql_GetById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 1
		rows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "net_weight", "expiration_rate",
			"recommended_freezing_temperature", "freezing_rate", "product_type_id",
			"seller_id", "width", "height", "length",
		}).AddRow(productId, "PROD001", "Product 1", 10.5, 30, -5.5, 10, 1, 1, 20.5, 10.5, 30.5)

		mock.ExpectQuery("SELECT (.+) FROM products WHERE id = \\?").WithArgs(productId).WillReturnRows(rows)

		repo := ProductSql{db: db}

		// Act
		product, err := repo.GetById(productId)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, product)
		require.Equal(t, productId, product.ID)
		require.Equal(t, "PROD001", *product.ProductCode)
		require.Equal(t, "Product 1", *product.Description)
		require.Equal(t, 10.5, *product.NetWeight)
		require.Equal(t, 30, *product.ExpirationRate)
		require.Equal(t, -5.5, *product.RecommendedFreezingTemperature)
		require.Equal(t, 10, *product.FreezingRate)
		require.Equal(t, 1, *product.ProductTypeId)
		require.Equal(t, 1, *product.SellerId)
		require.Equal(t, 20.5, *product.Width)
		require.Equal(t, 10.5, *product.Height)
		require.Equal(t, 30.5, *product.Length)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 999
		rows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "net_weight", "expiration_rate",
			"recommended_freezing_temperature", "freezing_rate", "product_type_id",
			"seller_id", "width", "height", "length",
		})
		mock.ExpectQuery("SELECT (.+) FROM products WHERE id = \\?").WithArgs(productId).WillReturnRows(rows)

		repo := ProductSql{db: db}

		// Act
		product, err := repo.GetById(productId)

		// Assert
		require.Error(t, err)
		require.Nil(t, product)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 1
		mock.ExpectQuery("SELECT (.+) FROM products WHERE id = \\?").WithArgs(productId).WillReturnError(sql.ErrConnDone)

		repo := ProductSql{db: db}

		// Act
		product, err := repo.GetById(productId)

		// Assert
		require.Error(t, err)
		require.Nil(t, product)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
func TestProductSql_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 1
		mock.ExpectExec("DELETE FROM products WHERE id = \\?").WithArgs(productId).WillReturnResult(sqlmock.NewResult(0, 1))

		repo := ProductSql{db: db}

		// Act
		err = repo.Delete(productId)

		// Assert
		require.NoError(t, err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 999
		mock.ExpectExec("DELETE FROM products WHERE id = \\?").WithArgs(productId).WillReturnResult(sqlmock.NewResult(0, 0))

		repo := ProductSql{db: db}

		// Act
		err = repo.Delete(productId)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 1
		mock.ExpectExec("DELETE FROM products WHERE id = \\?").WithArgs(productId).WillReturnError(errors.New("database error"))

		repo := ProductSql{db: db}

		// Act
		err = repo.Delete(productId)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict], err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows_affected_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 1
		result := sqlmock.NewErrorResult(errors.New("rows affected error"))

		mock.ExpectExec("DELETE FROM products WHERE id = \\?").WithArgs(productId).WillReturnResult(result)

		repo := ProductSql{db: db}

		// Act
		err = repo.Delete(productId)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
func TestProductSql_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		product := models.Product{
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    stringPtr("PROD001"),
				Description:                    stringPtr("Product 1"),
				NetWeight:                      float64Ptr(10.5),
				ExpirationRate:                 intPtr(30),
				RecommendedFreezingTemperature: float64Ptr(-5.5),
				FreezingRate:                   intPtr(10),
				ProductTypeId:                  intPtr(1),
				SellerId:                       intPtr(1),
			},
			Dimensions: models.Dimensions{
				Width:  float64Ptr(20.5),
				Height: float64Ptr(10.5),
				Length: float64Ptr(30.5),
			},
		}

		mock.ExpectExec("INSERT INTO products").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		repo := ProductSql{db: db}

		// Act
		createdProduct, err := repo.Create(product)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, createdProduct)
		require.Equal(t, 1, createdProduct.ID)
		require.Equal(t, *product.ProductCode, *createdProduct.ProductCode)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate_product_code", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		product := models.Product{
			ProductAttributes: models.ProductAttributes{
				ProductCode: stringPtr("PROD001"),
			},
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1062,
			Message: "Duplicate entry 'PROD001' for key 'product_code'",
		}

		mock.ExpectExec("INSERT INTO products").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
		).WillReturnError(mysqlErr)

		repo := ProductSql{db: db}

		// Act
		createdProduct, err := repo.Create(product)

		// Assert
		require.Error(t, err)
		require.Nil(t, createdProduct)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict].Code, serviceErr.Code)
		require.Equal(t, "product_code already exists", serviceErr.Message)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid_product_type_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		product := models.Product{
			ProductAttributes: models.ProductAttributes{
				ProductTypeId: intPtr(999), // Invalid product type ID
			},
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`frescos`.`products`, CONSTRAINT `fk_products_product_type_id` FOREIGN KEY (`product_type_id`) REFERENCES `product_types` ([id](cci:1://file:///Users/agaggino/Documents/Bootcamp/W17-G7-Bootcamp/pkg/models/product.go:66:0-98:1)))",
		}

		mock.ExpectExec("INSERT INTO products").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
		).WillReturnError(mysqlErr)

		repo := ProductSql{db: db}

		// Act
		createdProduct, err := repo.Create(product)

		// Assert
		require.Error(t, err)
		require.Nil(t, createdProduct)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)
		require.Equal(t, "invalid product_type_id", serviceErr.Message)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid_seller_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		product := models.Product{
			ProductAttributes: models.ProductAttributes{
				SellerId: intPtr(999), // Invalid seller ID
			},
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`frescos`.`products`, CONSTRAINT `fk_products_seller_id` FOREIGN KEY (`seller_id`) REFERENCES `sellers` ([id](cci:1://file:///Users/agaggino/Documents/Bootcamp/W17-G7-Bootcamp/pkg/models/product.go:66:0-98:1)))",
		}

		mock.ExpectExec("INSERT INTO products").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
		).WillReturnError(mysqlErr)

		repo := ProductSql{db: db}

		// Act
		createdProduct, err := repo.Create(product)

		// Assert
		require.Error(t, err)
		require.Nil(t, createdProduct)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)
		require.Equal(t, "invalid seller_id", serviceErr.Message)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("last_insert_id_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		product := models.Product{
			ProductAttributes: models.ProductAttributes{
				ProductCode: stringPtr("PROD001"),
			},
		}

		result := sqlmock.NewErrorResult(errors.New("last insert id error"))

		mock.ExpectExec("INSERT INTO products").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
		).WillReturnResult(result)

		repo := ProductSql{db: db}

		// Act
		createdProduct, err := repo.Create(product)

		// Assert
		require.Error(t, err)
		require.Nil(t, createdProduct)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("other_database_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		product := models.Product{
			ProductAttributes: models.ProductAttributes{
				ProductCode: stringPtr("PROD001"),
			},
		}

		mock.ExpectExec("INSERT INTO products").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
		).WillReturnError(errors.New("database error"))

		repo := ProductSql{db: db}

		// Act
		createdProduct, err := repo.Create(product)

		// Assert
		require.Error(t, err)
		require.Nil(t, createdProduct)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
func TestProductSql_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 1
		product := models.Product{
			ID: productId,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    stringPtr("PROD001"),
				Description:                    stringPtr("Updated Product"),
				NetWeight:                      float64Ptr(15.5),
				ExpirationRate:                 intPtr(45),
				RecommendedFreezingTemperature: float64Ptr(-10.5),
				FreezingRate:                   intPtr(15),
				ProductTypeId:                  intPtr(2),
				SellerId:                       intPtr(2),
			},
			Dimensions: models.Dimensions{
				Width:  float64Ptr(25.5),
				Height: float64Ptr(15.5),
				Length: float64Ptr(35.5),
			},
		}

		mock.ExpectExec("UPDATE products SET").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
			product.ID,
		).WillReturnResult(sqlmock.NewResult(0, 1))

		repo := ProductSql{db: db}

		// Act
		err = repo.Update(productId, product)

		// Assert
		require.NoError(t, err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 999
		product := models.Product{
			ID: productId,
		}

		mock.ExpectExec("UPDATE products SET").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
			product.ID,
		).WillReturnResult(sqlmock.NewResult(0, 0))

		repo := ProductSql{db: db}

		// Act
		err = repo.Update(productId, product)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate_product_code", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 1
		product := models.Product{
			ID: productId,
			ProductAttributes: models.ProductAttributes{
				ProductCode: stringPtr("PROD002"), // Already exists for another product
			},
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1062,
			Message: "Duplicate entry 'PROD002' for key 'product_code'",
		}

		mock.ExpectExec("UPDATE products SET").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
			product.ID,
		).WillReturnError(mysqlErr)

		repo := ProductSql{db: db}

		// Act
		err = repo.Update(productId, product)

		// Assert
		require.Error(t, err)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict].Code, serviceErr.Code)
		require.Equal(t, "product_code already exists", serviceErr.Message)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid_product_type_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 1
		product := models.Product{
			ID: productId,
			ProductAttributes: models.ProductAttributes{
				ProductTypeId: intPtr(999), // Invalid product type ID
			},
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`frescos`.`products`, CONSTRAINT `fk_products_product_type_id` FOREIGN KEY (`product_type_id`) REFERENCES `product_types` ([id](cci:1://file:///Users/agaggino/Documents/Bootcamp/W17-G7-Bootcamp/pkg/models/product.go:66:0-98:1)))",
		}

		mock.ExpectExec("UPDATE products SET").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
			product.ID,
		).WillReturnError(mysqlErr)

		repo := ProductSql{db: db}

		// Act
		err = repo.Update(productId, product)

		// Assert
		require.Error(t, err)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)
		require.Equal(t, "invalid product_type_id", serviceErr.Message)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid_seller_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := 1
		product := models.Product{
			ID: productId,
			ProductAttributes: models.ProductAttributes{
				SellerId: intPtr(999), // Invalid seller ID
			},
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`frescos`.`products`, CONSTRAINT `fk_products_seller_id` FOREIGN KEY (`seller_id`) REFERENCES `sellers` ([id](cci:1://file:///Users/agaggino/Documents/Bootcamp/W17-G7-Bootcamp/pkg/models/product.go:66:0-98:1)))",
		}

		mock.ExpectExec("UPDATE products SET").WithArgs(
			product.ProductCode,
			product.Description,
			product.NetWeight,
			product.ExpirationRate,
			product.RecommendedFreezingTemperature,
			product.FreezingRate,
			product.ProductTypeId,
			product.SellerId,
			product.Width,
			product.Height,
			product.Length,
			product.ID,
		).WillReturnError(mysqlErr)

		repo := ProductSql{db: db}

		// Act
		err = repo.Update(productId, product)

		// Assert
		require.Error(t, err)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)
		require.Equal(t, "invalid seller_id", serviceErr.Message)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestProductSql_GetProductRecords(t *testing.T) {
	t.Run("success_with_records", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := intPtr(1)
		rows := sqlmock.NewRows([]string{
			"id", "description", "records_count",
		}).AddRow(1, "Product 1", 5)

		mock.ExpectQuery("SELECT (.+) FROM products p INNER JOIN product_records pr").WithArgs(productId, productId).WillReturnRows(rows)

		repo := ProductSql{db: db}

		// Act
		records, err := repo.GetProductRecords(productId)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, records)
		require.Equal(t, 1, len(records))
		require.Equal(t, 1, records[0].ProductId)
		require.Equal(t, "Product 1", records[0].Description)
		require.Equal(t, 5, records[0].RecordsCount)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_with_no_records_but_product_exists", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := intPtr(1)

		// No records found in join query
		emptyRows := sqlmock.NewRows([]string{
			"id", "description", "records_count",
		})
		mock.ExpectQuery("SELECT (.+) FROM products p INNER JOIN product_records pr").WithArgs(productId, productId).WillReturnRows(emptyRows)

		// But product exists
		productRows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "net_weight", "expiration_rate",
			"recommended_freezing_temperature", "freezing_rate", "product_type_id",
			"seller_id", "width", "height", "length",
		}).AddRow(1, "PROD001", "Product 1", 10.5, 30, -5.5, 10, 1, 1, 20.5, 10.5, 30.5)

		mock.ExpectQuery("SELECT (.+) FROM products WHERE id = \\?").WithArgs(*productId).WillReturnRows(productRows)

		repo := ProductSql{db: db}

		// Act
		records, err := repo.GetProductRecords(productId)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, records)
		require.Equal(t, 1, len(records))
		require.Equal(t, 1, records[0].ProductId)
		require.Equal(t, "Product 1", records[0].Description)
		require.Equal(t, 0, records[0].RecordsCount)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := intPtr(1)
		mock.ExpectQuery("SELECT (.+) FROM products p INNER JOIN product_records pr").WithArgs(productId, productId).WillReturnError(sql.ErrConnDone)

		repo := ProductSql{db: db}

		// Act
		records, err := repo.GetProductRecords(productId)

		// Assert
		require.Error(t, err)
		require.Nil(t, records)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		productId := intPtr(1)
		rows := sqlmock.NewRows([]string{
			"id", "description", "records_count",
		}).AddRow("not_an_int", "Product 1", 5) // This will cause a scan error

		mock.ExpectQuery("SELECT (.+) FROM products p INNER JOIN product_records pr").WithArgs(productId, productId).WillReturnRows(rows)

		repo := ProductSql{db: db}

		// Act
		records, err := repo.GetProductRecords(productId)

		// Assert
		require.NoError(t, err)           // The function continues even if there's a scan error
		require.Equal(t, 0, len(records)) // But no records are added

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
