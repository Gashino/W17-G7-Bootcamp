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
)

// NewWarehouseDefault is a function that returns a new instance of VehicleDefault
func NewCarryDefault(sv service.CarryService) *CarryDefault {
	return &CarryDefault{sv: sv}
}

// Struct handler for Warehouses
type CarryDefault struct {
	// sv is the service that will be used by the handler
	sv service.CarryService
}

func (h *CarryDefault) SearchByLocality() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")

		// Inicializa idInt con -1
		idInt := -1

		// Si el parámetro id no está vacío, intenta convertirlo a un int
		if id != "" {
			var err error
			idInt, err = strconv.Atoi(id)
			if err != nil {
				response.Error(w, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, pkg.ServiceErrors[pkg.ErrNotFound].Error())
				return
			}
		}

		data, err := h.sv.SearchByLocality(idInt)

		if err != nil {
			fmt.Println(err.Error())
			response.Error(w, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, pkg.ServiceErrors[pkg.ErrNotFound].Error())
			return
		}

		// response
		writeResponse(w, http.StatusCreated, map[string]any{
			"data": data,
		}, nil)
	}
}

func (h *CarryDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Almaceno la del body
		var data models.CarryDoc
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

		carry := models.Carry{
			Cid:         data.Cid,
			CompanyName: data.CompanyName,
			Address:     data.Address,
			Telephone:   data.Telephone,
			LocalityId:  data.LocalityId,
		}

		// Llamo al service
		carry, err = h.sv.Create(carry)

		if err != nil {
			svcErr := pkg.ServiceErrors[pkg.ErrInternalServer]
			svcErr.InternalError = fmt.Errorf(err.Error())
			response.Error(w, svcErr.ResponseCode, svcErr.Error())
			return
		}

		writeResponse(w, http.StatusCreated, map[string]any{
			"data": data,
		}, nil)
		return
	}
}
