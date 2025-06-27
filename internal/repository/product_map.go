package repository

import (
	"app/pkg"
	"app/pkg/models"
)

type ProductMap struct {
	db map[int]models.Product
}

func (p ProductMap) GetById(id int) (*models.Product, error) {
	if value, exist := p.db[id]; exist {
		return &value, nil
	} else {
		return nil, pkg.ServiceErrors[pkg.ErrNotFound]
	}
}

func (p ProductMap) GetAll() map[int]models.Product {
	return p.db
}

func NewProductMap(db map[int]models.Product) *ProductMap {
	defaultDb := make(map[int]models.Product)
	if db != nil {
		defaultDb = db
	}
	return &ProductMap{db: defaultDb}
}
