package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
)

// NewWarehouseDefault is a function that returns a new instance of VehicleDefault
func NewWarehouseDefault(sv service.WarehouseService) *WarehouseDefault {
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
			response.Error(w, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		// response
		data := make(map[int]models.WarehouseDoc)
		for key, value := range v {
			data[key] = models.WarehouseDoc{
				ID:             value.ID,
				WarehouseCode:  value.WarehouseCode,
				Address:        value.Address,
				Telephone:      value.Telephone,
				MinCapacity:    value.MinCapacity,
				MinTemperature: value.MinTemperature,
			}
		}

		writeResponse(w, http.StatusOK, map[string]any{
			"data": data,
		}, nil)

	}
}

func (h *WarehouseDefault) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		idInt, err := strconv.Atoi(id)

		if err != nil {
			response.Error(w, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, pkg.ServiceErrors[pkg.ErrNotFound].Error())
			return
		}

		value, err := h.sv.FindByID(idInt)

		if err != nil {

			response.Error(w, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, pkg.ServiceErrors[pkg.ErrNotFound].Error())
			return
		}

		// response

		data := models.Warehouse{
			ID:             value.ID,
			WarehouseCode:  value.WarehouseCode,
			Address:        value.Address,
			Telephone:      value.Telephone,
			MinCapacity:    value.MinCapacity,
			MinTemperature: value.MinTemperature,
		}

		writeResponse(w, http.StatusCreated, map[string]any{
			"data": data,
		}, nil)
	}
}

func (h *WarehouseDefault) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Almaceno la del body
		var data models.WarehouseDoc
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		defer r.Body.Close()

		// Si no se puede hacer cancelo
		if err != nil {
			response.Error(w, pkg.ServiceErrors[pkg.ErrUnprocessableEntity].ResponseCode, pkg.ServiceErrors[pkg.ErrUnprocessableEntity].Error())
			return
		}

		// Si algun campo esta vacio cancelo
		valid := data.AreFieldsValid()
		if !valid {

			response.Error(w, pkg.ServiceErrors[pkg.ErrUnprocessableEntity].ResponseCode, pkg.ServiceErrors[pkg.ErrUnprocessableEntity].Error())
			return
		}

		// Llamo al service
		warehouse, err := h.sv.Add(data)

		if err != nil {
			svcErr := pkg.ServiceErrors[pkg.ErrInternalServer]
			svcErr.InternalError = fmt.Errorf(err.Error())
			response.Error(w, svcErr.ResponseCode, svcErr.Error())
			return
		}

		warehouse_doc := models.WarehouseDoc{
			ID:             warehouse.ID,
			WarehouseCode:  warehouse.WarehouseCode,
			Address:        warehouse.Address,
			Telephone:      warehouse.Telephone,
			MinCapacity:    warehouse.MinCapacity,
			MinTemperature: warehouse.MinTemperature,
		}

		writeResponse(w, http.StatusCreated, map[string]any{
			"data": warehouse_doc,
		}, nil)
		return
	}
}

func (h *WarehouseDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)

		// Almaceno la del body
		var data models.WarehouseDoc
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&data)
		defer r.Body.Close()

		// Si no se puede hacer cancelo
		if err != nil {
			response.Error(w, pkg.ServiceErrors[pkg.ErrUnprocessableEntity].ResponseCode, pkg.ServiceErrors[pkg.ErrUnprocessableEntity].Error())
			return
		}

		// Llamo al service
		warehouse, err := h.sv.Update(idInt, data)

		if err != nil {
			svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
			svcErr.InternalError = fmt.Errorf(err.Error())
			response.Error(w, svcErr.ResponseCode, svcErr.Error())
			return
		}

		warehouse_doc := models.WarehouseDoc{
			ID:             warehouse.ID,
			WarehouseCode:  warehouse.WarehouseCode,
			Address:        warehouse.Address,
			Telephone:      warehouse.Telephone,
			MinCapacity:    warehouse.MinCapacity,
			MinTemperature: warehouse.MinTemperature,
		}

		//response.JSON(w, http.StatusCreated, warehouse_doc)
		writeResponse(w, http.StatusCreated, map[string]any{
			"data": warehouse_doc,
		}, nil)
		return
	}
}

func (h *WarehouseDefault) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.Error(w, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, pkg.ServiceErrors[pkg.ErrNotFound].Error())
			return
		}
		err = h.sv.Delete(idInt)

		if err != nil {
			svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
			svcErr.InternalError = fmt.Errorf(err.Error())
			response.Error(w, svcErr.ResponseCode, svcErr.Error())
			return
		}

		response.JSON(w, http.StatusNoContent, map[string]any{
			"message": "success",
			"data":    "Deleted",
		})
		return

	}

}
