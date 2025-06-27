package repository

import "app/pkg/models"

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
