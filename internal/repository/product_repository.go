package repository

import "app/pkg/models"

type ProductRepository interface {
	GetAll() map[int]models.Product
	GetById(id int) (*models.Product, error)
	Delete(id int) error
	Create(product models.Product) error
	Update(id int, product models.Product) error
}
