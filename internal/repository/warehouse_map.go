package repository

import (
	"app/pkg/models"
	"errors"
)

// NewVehicleMap is a function that returns a new instance of VehicleMap
func NewWarehouseMap(sectionDb *map[int]models.Section, employeeDb *map[int]models.Employee, warehouseDb *map[int]models.Warehouse) *WarehouseMap {
	return &WarehouseMap{db: warehouseDb, sectionDb: sectionDb, employeeDb: employeeDb}
}

// Struct for Warehouse Repository
type WarehouseMap struct {
	// db is a map of warehouse
	db         *map[int]models.Warehouse
	sectionDb  *map[int]models.Section
	employeeDb *map[int]models.Employee
}

func (r *WarehouseMap) FindAll() (v map[int]models.Warehouse, err error) {
	v = make(map[int]models.Warehouse)

	// copy db
	for key, value := range *r.db {
		v[key] = value
	}

	return
}

func (r *WarehouseMap) FindOne(id int) (v models.Warehouse, err error) {
	for _, value := range *r.db {
		if value.ID == id {
			return value, err
		}

	}
	err = errors.New("warehouse not found")
	return
}

func (r *WarehouseMap) Add(v models.Warehouse) (err error) {
	(*r.db)[v.ID] = v
	return
}

func (r *WarehouseMap) FindWarehouseByCode(code string) (v models.Warehouse, err error) {
	for _, value := range *r.db {
		if value.WarehouseCode == code {
			return value, err
		}

	}
	err = errors.New("warehouse not found")
	return
}

func (r *WarehouseMap) FindAvailableID() (id int, err error) {
	if len(*r.db) == 0 {
		// Si el mapa está vacío
		id = 1
		return
	}

	maxID := 0
	for id := range *r.db {
		if id > maxID {
			maxID = id
		}
	}

	id = maxID + 1
	return
}

func (r *WarehouseMap) Delete(id int) (err error) {

	// Verifico Section
	for _, s := range *r.sectionDb {
		if s.WarehouseID == id {
			err = errors.New("FK restriction with Section")
			return
		}
	}

	// Verifico Employee
	for _, e := range *r.employeeDb {
		if e.WarehouseID == id {
			err = errors.New("FK restriction with Employee")
			return
		}
	}

	delete(*r.db, id)
	return
}
