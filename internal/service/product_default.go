package service

import "app/internal/repository"

func NewProductDefault(rp repository.ProductRepository) *ProductService {
	return &ProductService{rp: rp}
}

type ProductService struct {
	rp repository.ProductRepository
}
