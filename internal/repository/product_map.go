package repository

import "app/pkg/models"

type ProductMap struct {
	db map[int]models.Product
}

func NewProductMap(db map[int]models.Product) *ProductMap {
	defaultDb := make(map[int]models.Product)
	if db != nil {
		defaultDb = db
	}
	return &ProductMap{db: defaultDb}
}
