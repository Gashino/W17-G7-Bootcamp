package repository

import (
	"app/pkg"
	"app/pkg/models"
	"net/http"
	"strings"
)

type ProductMap struct {
	db            *map[int]models.Product
	dbProductType *map[int]models.ProductType
	lastId        int
}

func (p *ProductMap) Update(id int, product models.Product) error {

	if err := p.validateProductCode(product); err != nil {
		return err
	}
	(*p.db)[id] = product

	return nil
}

func (p *ProductMap) Create(product models.Product) error {
	if err := p.validateProductCode(product); err != nil {
		return err
	}

	product.ID = p.lastId
	(*p.db)[product.ID] = product
	p.lastId++

	return nil
}

func (p *ProductMap) Delete(id int) error {
	if _, exist := (*p.db)[id]; exist {
		delete(*p.db, id)
	} else {
		return pkg.ServiceErrors[pkg.ErrNotFound]
	}

	return nil
}

func (p *ProductMap) GetById(id int) (*models.Product, error) {
	if value, exist := (*p.db)[id]; exist {
		return &value, nil
	} else {
		return nil, pkg.ServiceErrors[pkg.ErrNotFound]
	}
}

func (p *ProductMap) GetAll() map[int]models.Product {
	return *p.db
}

func NewProductMap(dbProduct *map[int]models.Product, dbProductType *map[int]models.ProductType) *ProductMap {
	return &ProductMap{db: dbProduct, dbProductType: dbProductType, lastId: len(*dbProduct) + 1}
}

func (p *ProductMap) validateProductCode(product models.Product) error {
	for _, value := range *p.db {
		if value.ID == product.ID {
			continue
		}

		if strings.ToLower(value.ProductCode) == strings.ToLower(product.ProductCode) {
			return pkg.ServiceError{
				Code:         0,
				ResponseCode: http.StatusConflict,
				Message:      "Product code already exist",
			}
		}
	}
	return nil
}
