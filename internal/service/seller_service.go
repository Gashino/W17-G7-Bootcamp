package service

import "app/pkg/models"

// SellerService is an interface that represents a Seller service
type SellerService interface {
	// FindAll is a method that returns a map of all Sellers
	FindAll() (v map[int]models.Seller, err error)
	Create(seller models.Seller) error
	GetById(id int) (models.Seller, error)
	UpdateFields(id int, data models.SellerCreateRequest) (models.Seller, error)
	DeleteSeller(id int) error
}
