package repository

import "app/pkg/models"

// LocalityRepository is an interface that represents a Locality repository
type LocalityRepository interface {
	// GetById is a method that returns a Locality by its ID
	GetById(id int) (models.Locality, error)
	// Create is a method that creates a new Locality
	Create(locality models.Locality) (models.Locality, error)
	// GetCantSellersByLocality is a method that returns sellers count by locality
	GetCantSellersByLocality(id int) (models.LocalityBySellerResponse, error)
}
