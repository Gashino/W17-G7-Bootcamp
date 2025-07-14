package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

func NewCarryDefault(rp repository.CarryRepository) *CarryDefault {
	return &CarryDefault{rp: rp}
}

// Struct for Warehouse Service
type CarryDefault struct {
	// rp is the repository that will be used by the service
	rp repository.CarryRepository
}

func (s *CarryDefault) SearchByLocality(locality_id int) (c map[int]models.Carry, err error) {
	c, err = s.rp.SearchByLocality(locality_id)
	return
}

func (s *CarryDefault) Add(v models.Carry) (w models.Carry, err error) {
	w = v
	w, err = s.rp.Create(w)
	return
}
