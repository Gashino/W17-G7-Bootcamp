package service

import "app/pkg/models"

// LocalityService is an interface that represents a Locality service
type LocalityService interface {
	Create(seller models.Locality) (models.Locality, error)
	GetById(id int) (models.Locality, error)
	GetCantSellersByLocality(id int) (models.LocalityBySellerResponse, error)
}
