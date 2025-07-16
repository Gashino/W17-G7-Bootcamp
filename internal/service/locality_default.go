package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

// NewLocalityDefault is a function that returns a new instance of LocalityDefault
func NewLocalityDefault(rp repository.LocalityRepository) *LocalityDefault {
	return &LocalityDefault{rp: rp}
}

// LocalityDefault is a struct that represents the default service for Localities
type LocalityDefault struct {
	// rp is the repository that will be used by the service
	rp repository.LocalityRepository
}

// Create is a method that creates a new Locality
func (s *LocalityDefault) Create(locality models.Locality) (models.Locality, error) {
	return s.rp.Create(locality)
}

// GetById is a method that returns a Locality by its ID
func (s *LocalityDefault) GetById(id int) (models.Locality, error) {
	return s.rp.GetById(id)
}

// GetCantSellersByLocality is a method that returns sellers count by locality
func (s *LocalityDefault) GetCantSellersByLocality(id int) (models.LocalityBySellerResponse, error) {
	return s.rp.GetCantSellersByLocality(id)
}
