package service

import "app/pkg/models"

// LocalityService is an interface that represents a Locality service
type LocalityService interface {
	// Create is a method that creates a new Locality
	Create(locality models.Locality) (models.Locality, error)
	// GetById is a method that returns a Locality by its ID
	GetById(id int) (models.Locality, error)
	// GetCantSellersByLocality is a method that returns sellers count by locality
	GetCantSellersByLocality(id int) (models.LocalityBySellerResponse, error)
}
