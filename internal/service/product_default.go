package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

func NewProductDefault(rp repository.ProductRepository) *ProductService {
	return &ProductService{rp: rp}
}

type ProductService struct {
	rp repository.ProductRepository
}

func (p ProductService) Create(product models.Product) error {
	return p.rp.Create(product)
}

func (p ProductService) Delete(id int) error {
	return p.rp.Delete(id)
}

func (p ProductService) GetById(id int) (*models.Product, error) {
	return p.rp.GetById(id)
}

func (p ProductService) GetAll() (map[int]models.Product, error) {
	result := p.rp.GetAll()
	return result, nil
}
