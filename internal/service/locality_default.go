package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

// NewLocalityDefault is a function that returns a new instance of LocalityDefault
func NewLocalityDefault(rp repository.LocalityRepository) *LocalityDefault {
	return &LocalityDefault{rp: rp}
}

// LocalityDefault is a struct that represents the default service for Sellers
type LocalityDefault struct {
	// rp is the repository that will be used by the service
	rp repository.LocalityRepository
}

func (s *LocalityDefault) Create(seller models.Locality) (models.Locality, error) {
	return s.rp.Create(seller)
}

func (s *LocalityDefault) GetById(id int) (models.Locality, error) {
	return s.rp.GetById(id)
}

func (s *LocalityDefault) GetCantSellersByLocality(id int) (models.LocalityBySellerResponse, error) {
	return s.rp.GetCantSellersByLocality(id)
}
