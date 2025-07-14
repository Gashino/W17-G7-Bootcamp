package repository

import (
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

func (p ProductSql) GetAll() map[int]models.Product {
	//TODO implement me
	panic("implement me")
}

func (p ProductSql) GetById(id int) (*models.Product, error) {
	//TODO implement me
	panic("implement me")
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
