package handler

import (
	"app/internal/service"
	"app/pkg/models"
	"net/http"

	"github.com/bootcamp-go/web/response"
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
		v, err := h.sv.FindAll()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		// response
		data := make(map[int]models.WarehouseDoc)
		for key, value := range v {
			data[key] = models.WarehouseDoc{
				Id:             value.Id,
				WarehouseCode:  value.WarehouseCode,
				Address:        value.Address,
				Telephone:      value.Telephone,
				MinCapacity:    value.MinCapacity,
				MinTemperature: value.MinTemperature,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}
