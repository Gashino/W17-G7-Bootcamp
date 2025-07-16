package service

import "app/pkg/models"

type CarryService interface {
	Create(v models.Carry) (c models.Carry, err error)
	SearchByLocality(locality_id int) (v map[string]models.CarryByLocality, err error)
}
