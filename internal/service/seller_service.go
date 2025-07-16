package service

import "app/pkg/models"

// SellerService is an interface that represents a Seller service
type SellerService interface {
	// FindAll is a method that returns a map of all Sellers
	FindAll() (v map[int]models.Seller, err error)
	// Create is a method that creates a new Seller
	Create(seller models.Seller) (models.Seller, error)
	// GetById is a method that returns a Seller by its ID
	GetById(id int) (models.Seller, error)
	// UpdateFields is a method that updates a Seller's fields
	UpdateFields(id int, data models.SellerCreateRequest) (models.Seller, error)
	// DeleteSeller is a method that deletes a Seller
	DeleteSeller(id int) error
}
