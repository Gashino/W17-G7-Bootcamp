package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"strings"
)

type ProductSql struct {
	db *sql.DB
}

func NewProductSqlRepository(db *sql.DB) ProductRepository {

	return &ProductSql{
		db: db,
	}
}

const (
	GetAll = `SELECT 
		id,
		product_code,
		description,
		net_weight,
		expiration_rate,
		recommended_freezing_temperature,
		freezing_rate,
		product_type_id,
		seller_id,
		width,
		height,
		length FROM products`
	GetById = `SELECT
		id,
		product_code,
		description,
		net_weight,
		expiration_rate,
		recommended_freezing_temperature,
		freezing_rate,
		product_type_id,
		seller_id,
		width,
		height,
		length FROM products WHERE id = ?`
	Delete = `DELETE FROM products WHERE id = ?`
	Create = `INSERT INTO products (
		product_code,
		description,
		net_weight,
		expiration_rate,
		recommended_freezing_temperature,
		freezing_rate,
		product_type_id,
		seller_id,
		width,
		height,
		length
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

func (p ProductSql) GetAll() map[int]models.Product {
	productMap := make(map[int]models.Product)
	rows, err := p.db.Query(GetAll)

	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var prod models.Product
		err := rows.Scan(
			&prod.ID,
			&prod.ProductCode,
			&prod.Description,
			&prod.NetWeight,
			&prod.ExpirationRate,
			&prod.RecommendedFreezingTemperature,
			&prod.FreezingRate,
			&prod.ProductTypeId,
			&prod.SellerId,
			&prod.Width,
			&prod.Height,
			&prod.Length,
		)
		if err != nil {
			continue
		}
		productMap[prod.ID] = prod
	}

	return productMap
}

func (p ProductSql) GetById(id int) (*models.Product, error) {
	row := p.db.QueryRow(GetById, id)

	if err := row.Err(); err != nil {
		return nil, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	var prod models.Product

	if err := row.Scan(&prod.ID,
		&prod.ProductCode,
		&prod.Description,
		&prod.NetWeight,
		&prod.ExpirationRate,
		&prod.RecommendedFreezingTemperature,
		&prod.FreezingRate,
		&prod.ProductTypeId,
		&prod.SellerId,
		&prod.Width,
		&prod.Height,
		&prod.Length); err != nil {
		return nil, pkg.ServiceErrors[pkg.ErrNotFound]
	}

	return &prod, nil

}

func (p ProductSql) Delete(id int) error {
	rows, err := p.db.Exec(Delete, id)

	if err != nil {
		return pkg.ServiceErrors[pkg.ErrConflict]
	}

	rAffected, errRows := rows.RowsAffected()

	if errRows != nil {
		return pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	if rAffected == 0 {
		return pkg.ServiceErrors[pkg.ErrNotFound]
	}

	return nil
}

func (p ProductSql) Create(product models.Product) (*models.Product, error) {
	result, err := p.db.Exec(Create,
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
	)
	if err != nil {
		var errorResponse pkg.ServiceError
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				errorResponse = pkg.ServiceErrors[pkg.ErrConflict]
				errorResponse.Message = "product_code already exists"
				return nil, errorResponse
			case 1452:
				errorResponse = pkg.ServiceErrors[pkg.ErrNotFound]
				if strings.Contains(mysqlErr.Message, "product_type_id") {
					errorResponse.Message = "invalid product_type_id"
				} else if strings.Contains(mysqlErr.Message, "seller_id") {
					errorResponse.Message = "invalid seller_id"
				}
				return nil, errorResponse
			}
		}
		return nil, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	lastId, errId := result.LastInsertId()
	if errId != nil {
		return nil, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	product.ID = int(lastId)
	return &product, nil

}

func (p ProductSql) Update(id int, product models.Product) error {
	//TODO implement me
	panic("implement me")
}

//func (p *ProductSql) validateProductCode(product models.Product) error {
//	for _, value := range *p.db {
//		if value.ID == product.ID {
//			continue
//		}
//
//		if strings.ToLower(*value.ProductCode) == strings.ToLower(*product.ProductCode) {
//			return pkg.ServiceError{
//				Code:         0,
//				ResponseCode: http.StatusConflict,
//				Message:      "Product code already exist",
//			}
//		}
//	}
//	return nil
//}
//
