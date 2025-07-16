package repository

import "app/pkg/models"

// interface that represents a warehouse repository
type CarryRepository interface {
	Create(v models.Carry) (c models.Carry, err error)
	SearchByLocality(locality_id int) (v map[string]models.CarryByLocality, err error)
}
