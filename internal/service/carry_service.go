package service

import "app/pkg/models"

type CarryService interface {
	Create(v models.Carry) (c models.Carry, err error)
	SearchByLocality(locality_id int) (c map[int]models.Carry, err error)
}
