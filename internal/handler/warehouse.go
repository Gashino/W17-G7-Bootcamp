package handler

import (
	"app/internal/service"
	"net/http"
)

// NewVehicleDefault is a function that returns a new instance of VehicleDefault
func NewVehicleDefault(sv service.WarehouseService) *WarehouseDefault {
	return &WarehouseDefault{sv: sv}
}

// Struct handler for Warehouses
type WarehouseDefault struct {
	// sv is the service that will be used by the handler
	sv service.WarehouseService
}

func (h *WarehouseDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
