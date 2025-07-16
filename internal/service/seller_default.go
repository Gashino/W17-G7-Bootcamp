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

// Create is a method that creates a new Seller
func (s *SellerDefault) Create(seller models.Seller) (models.Seller, error) {
	return s.rp.Create(seller)
}

// GetById is a method that returns a Seller by its ID
func (s *SellerDefault) GetById(id int) (models.Seller, error) {
	return s.rp.GetById(id)
}

// UpdateFields is a method that updates a Seller's fields
func (s *SellerDefault) UpdateFields(id int, data models.SellerCreateRequest) (models.Seller, error) {
	return s.rp.UpdateFields(id, data)
}

// DeleteSeller is a method that deletes a Seller
func (s *SellerDefault) DeleteSeller(id int) error {
	return s.rp.DeleteSeller(id)
}
