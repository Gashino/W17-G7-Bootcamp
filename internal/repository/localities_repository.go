package repository

import "app/pkg/models"

// LocalityRepository is an interface that represents a Seller repository
type LocalityRepository interface {
	// FindAll is a method that returns a map of all Sellers
	GetById(id int) (models.Locality, error)
	Create(seller models.Locality) (models.Locality, error)
}
