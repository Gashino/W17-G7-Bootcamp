package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
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
		length FROM products WHERE ID = ?`
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
	//TODO implement me
	panic("implement me")
}

func (p ProductSql) Create(product models.Product) (*models.Product, error) {
	//TODO implement me
	panic("implement me")
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
//func (p *ProductSql) existsProductID(id int) bool {
//	for _, section := range *p.dbSection {
//		for _, batch := range section.ProductBatches {
//			for _, productID := range batch.ProductIDs {
//				if productID == id {
//					return true
//				}
//			}
//		}
//	}
//	return false
//}
