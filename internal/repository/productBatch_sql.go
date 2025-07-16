package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

func NewProductBatchSqlRepository(db *sql.DB) ProductBatchRepository {

	return &ProductBatchRepositorySql{
		db: db,
	}
}

type ProductBatchRepositorySql struct {
	db *sql.DB
}

func (r *ProductBatchRepositorySql) InsertProductBatch(pb models.ProductBatch) (productBatch models.ProductBatch, err error) {
	result, err := r.db.Exec(`
		INSERT INTO product_batches (
			batch_number,
			current_quantity,
			current_temperature,
			due_date,
			manufacturing_date,
			manufacturing_hour,
			minimum_temperature,
			initial_quantity,
			product_id,
			section_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		pb.BatchNumber,
		pb.CurrentQuantity,
		pb.CurrentTemperature,
		pb.DueDate,
		pb.ManufacturingDate,
		pb.ManufacturingHour,
		pb.MinumumTemperature,
		pb.InitialQuantity,
		pb.ProductId,
		pb.SectionId,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				svcErr := pkg.ServiceErrors[pkg.ErrConflict]
				svcErr.InternalError = fmt.Errorf("batch_number duplicado")
				return models.ProductBatch{}, svcErr
			case 1452:
				svcErr := pkg.ServiceErrors[pkg.ErrConflict]
				svcErr.InternalError = fmt.Errorf("product_id o section_id no existen")
				return models.ProductBatch{}, svcErr
			}
		}
		return models.ProductBatch{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	id, err := result.LastInsertId()
	if err != nil {
		return models.ProductBatch{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	pb.ID = int(id)
	return pb, nil
}
