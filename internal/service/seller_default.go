package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

// NewSellerDefault is a function that returns a new instance of SellerDefault
func NewSellerDefault(rp repository.SellerRepository) *SellerDefault {
	return &SellerDefault{rp: rp}
}

// SellerDefault is a struct that represents the default service for Sellers
type SellerDefault struct {
	// rp is the repository that will be used by the service
	rp repository.SellerRepository
}

// FindAll is a method that returns a map of all Sellers
func (s *SellerDefault) FindAll() (v map[int]models.Seller, err error) {
	v, err = s.rp.FindAll()
	return
}

func (s *SellerDefault) Create(seller models.Seller) error {
	return s.rp.Create(seller)
}

func (s *SellerDefault) GetById(id int) (models.Seller, error) {
	return s.rp.GetById(id)
}

func (s *SellerDefault) UpdateFields(id int, data models.SellerCreateRequest) (models.Seller, error) {
	return s.rp.UpdateFields(id, data)
}

func (s *SellerDefault) DeleteSeller(id int) error {
	return s.rp.DeleteSeller(id)
}
