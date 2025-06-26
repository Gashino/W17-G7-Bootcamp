package service

import "app/internal/repository"

func NewVehicleDefault(rp repository.WarehouseRepository) *WarehouseDefault {
	return &WarehouseDefault{rp: rp}
}

// Struct for Warehouse Service
type WarehouseDefault struct {
	// rp is the repository that will be used by the service
	rp repository.WarehouseRepository
}
