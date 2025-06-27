package repository

import (
	"app/pkg/models"
	"errors"
)

// NewVehicleMap is a function that returns a new instance of VehicleMap
func NewWarehouseMap(db map[int]models.Warehouse) *WarehouseMap {
	// default db
	defaultDb := make(map[int]models.Warehouse)
	if db != nil {
		defaultDb = db
	}
	return &WarehouseMap{db: defaultDb}
}

// Struct for Warehouse Repository
type WarehouseMap struct {
	// db is a map of warehouse
	db map[int]models.Warehouse
}

func (r *WarehouseMap) FindAll() (v map[int]models.Warehouse, err error) {
	v = make(map[int]models.Warehouse)

	// copy db
	for key, value := range r.db {
		v[key] = value
	}

	return
}

func (r *WarehouseMap) FindOne(id int) (v models.Warehouse, err error) {
	for _, value := range r.db {
		if value.Id == id {
			return value, err
		}

	}
	err = errors.New("warehouse not found")
	return
}
