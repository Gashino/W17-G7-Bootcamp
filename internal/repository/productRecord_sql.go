package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
)

const (
	Insert = `INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?,?,?,?)`
)

type ProductRecordSql struct {
	db *sql.DB
}

func NewProductRecordSqlRepository(db *sql.DB) ProductRecordRepository {
	return &ProductRecordSql{
		db: db,
	}
}

func (p ProductRecordSql) Insert(record models.ProductRecord) (*models.ProductRecord, error) {
	result, err := p.db.Exec(Insert,
		record.LastUpdateDate,
		record.PurchasePrice,
		record.SalePrice,
		record.ProductId,
	)

	if err != nil {
		var errorResponse pkg.ServiceError
		errorResponse = pkg.ServiceErrors[pkg.ErrNotFound]
		errorResponse.Message = "invalid product_id"
		return nil, errorResponse
	}

	lastId, errId := result.LastInsertId()
	if errId != nil {
		return nil, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	record.ID = int(lastId)
	return &record, nil

}
