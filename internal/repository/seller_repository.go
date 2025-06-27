package repository

import "app/pkg/models"

// SellerRepository is an interface that represents a Seller repository
type SellerRepository interface {
	// FindAll is a method that returns a map of all Sellers
	FindAll() (v map[int]models.Seller, err error)
	GetById(id int) (models.Seller, error)
	Create(seller models.Seller) error
	UpdateFields(id int, data models.SellerCreateRequest) (models.Seller, error)
	DeleteSeller(id int) error
}
