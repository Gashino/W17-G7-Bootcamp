package service

import (
	"app/internal/repository"
	"app/pkg/models"
	"errors"
)

func NewVehicleDefault(rp repository.WarehouseRepository) *WarehouseDefault {
	return &WarehouseDefault{rp: rp}
}

// Struct for Warehouse Service
type WarehouseDefault struct {
	// rp is the repository that will be used by the service
	rp repository.WarehouseRepository
}

func (s *WarehouseDefault) FindAll() (v map[int]models.Warehouse, err error) {
	v, err = s.rp.FindAll()
	return

}

func (s *WarehouseDefault) FindOne(id int) (v models.Warehouse, err error) {
	v, err = s.rp.FindOne(id)
	return
}

func (s *WarehouseDefault) Add(v models.WarehouseDoc) (err error) {
	// Verfico que el codigo sea unico
	_, err = s.rp.FindWarehouseByCode(v.WarehouseCode)

	if err == nil {
		err = errors.New("Warehouse Code is not unique")
		return
	}

	// Obtengo el nuevo ID
	id, err := s.rp.FindAvailableID()
	if err != nil {
		return
	}

	warehouse := models.Warehouse{
		Id:             id,
		WarehouseCode:  v.WarehouseCode,
		Address:        v.Address,
		Telephone:      v.Telephone,
		MinCapacity:    v.MinCapacity,
		MinTemperature: v.MinTemperature,
	}
	err = s.rp.Add(warehouse)

	return
}
