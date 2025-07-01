package service

import "app/pkg/models"

type IProductService interface {
	GetAll() (map[int]models.Product, error)
	GetById(id int) (*models.Product, error)
	Delete(id int) error
	Create(product models.Product) error
	Update(id int, product models.ProductDoc) (*models.Product, error)
}
