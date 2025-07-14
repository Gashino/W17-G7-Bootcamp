package service

import (
	"app/internal/repository"
	"app/pkg/models"
	"errors"
	"fmt"
)

func NewWarehouseDefault(rp repository.WarehouseRepository) *WarehouseDefault {
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

func (s *WarehouseDefault) FindByID(id int) (v models.Warehouse, err error) {
	v, err = s.rp.FindByID(id)
	return
}

func (s *WarehouseDefault) Add(v models.WarehouseDoc) (w models.Warehouse, err error) {
	// Verfico que el codigo sea unico
	_, err = s.rp.FindWarehouseByCode(v.WarehouseCode)

	if err == nil {
		err = errors.New("Warehouse Code is not unique")
		return
	}

	// Obtengo el nuevo ID
	//id, err := s.rp.FindAvailableID()
	//if err != nil {
	//	return
	//}

	w = models.Warehouse{
		WarehouseCode:  v.WarehouseCode,
		Address:        v.Address,
		Telephone:      v.Telephone,
		MinCapacity:    v.MinCapacity,
		MinTemperature: v.MinTemperature,
	}
	w, err = s.rp.Add(w)

	return
}

func (s *WarehouseDefault) Update(id int, v models.WarehouseDoc) (w models.Warehouse, err error) {
	// Obtengo el elemento
	warehouse, err := s.rp.FindByID(id)
	if err != nil {
		return
	}
	total_changes := 0
	fmt.Println(v.WarehouseCode, "-", warehouse.WarehouseCode)
	if v.WarehouseCode != "" && warehouse.WarehouseCode != v.WarehouseCode {
		// Verfico que el codigo sea unico
		_, err = s.rp.FindWarehouseByCode(v.WarehouseCode)
		if err == nil {
			err = errors.New("Warehouse Code is not unique")
			return
		}
		warehouse.WarehouseCode = v.WarehouseCode
		total_changes++
	}

	if v.Address != "" && warehouse.Address != v.Address {
		warehouse.Address = v.Address
		total_changes++
	}

	if v.Telephone != "" && warehouse.Telephone != v.Telephone {
		warehouse.Telephone = v.Telephone
		total_changes++
	}

	if v.MinCapacity != 0 && warehouse.MinCapacity != v.MinCapacity {
		warehouse.MinCapacity = v.MinCapacity
		total_changes++
	}

	if v.MinTemperature != 0 && warehouse.MinTemperature != v.MinTemperature {
		warehouse.MinTemperature = v.MinTemperature
		total_changes++
	}

	if total_changes > 0 {
		warehouse, err = s.rp.Add(warehouse)
	}

	w = warehouse
	return
}

func (s *WarehouseDefault) Delete(id int) (err error) {
	// Obtengo el elemento
	_, err = s.rp.FindByID(id)
	if err != nil {
		return
	}

	err = s.rp.Delete(id)
	return
}
