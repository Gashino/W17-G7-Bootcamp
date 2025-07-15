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

func (p ProductService) GetProductRecords(id *int) ([]models.ProductRecordResponse, error) {
	return p.rp.GetProductRecords(id)
}

func (p ProductService) Update(id int, product models.Product) (*models.Product, error) {
	productToUpdate, errId := p.rp.GetById(id)
	if errId != nil {
		return nil, errId
	}
	product.MapDataToStruct(productToUpdate)

	err := p.rp.Update(id, *productToUpdate)
	if err != nil {
		return nil, err
	}

	return productToUpdate, nil
}

func (p ProductService) Create(product models.Product) (*models.Product, error) {
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
