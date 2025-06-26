package service

import (
	"app/internal/repository"
	"app/pkg"
	"app/pkg/models"
	"fmt"
)

// NewVehicleDefault is a function that returns a new instance of VehicleDefault
func NewVehicleDefault(rp repository.VehicleRepository) *VehicleDefault {
	return &VehicleDefault{rp: rp}
}

// VehicleDefault is a struct that represents the default service for vehicles
type VehicleDefault struct {
	// rp is the repository that will be used by the service
	rp repository.VehicleRepository
}

// FindAll is a method that returns a map of all vehicles
func (s *VehicleDefault) FindAll() (v map[int]models.Vehicle, err error) {
	v, err = s.rp.FindAll()
	return
}

func (s *VehicleDefault) UpdateSpeedForId(id int, velocidad float64) (v models.Vehicle, err error) {
	err = s.ValidateSpeed(velocidad)
	if err != nil {
		return models.Vehicle{}, err
	}

	v, err = s.rp.UpdateSpeedForId(id, velocidad)
	if err != nil {
		return models.Vehicle{}, err
	}

	return v, nil
}

func (s *VehicleDefault) ValidateSpeed(velocidad float64) error {
	if velocidad <= 0.0 {
		svcErr := pkg.ServiceErrors[pkg.ErrBadRequest]
		svcErr.InternalError = fmt.Errorf("Velocidad mal formada o fuera de rango.")
		return svcErr
	}
	return nil
}
