package repository

import "app/pkg/models"

type ProductRepository interface {
	GetAll() map[int]models.Product
	GetById(id int) (*models.Product, error)
}
