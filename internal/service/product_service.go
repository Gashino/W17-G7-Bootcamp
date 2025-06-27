package service

import "app/pkg/models"

type IProductService interface {
	GetAll() (map[int]models.Product, error)
	GetById(id int) (*models.Product, error)
}
