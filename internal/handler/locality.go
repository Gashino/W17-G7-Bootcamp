package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
)

// NewLocalityDefault is a function that returns a new instance of LocalityDefault
func NewLocalityDefault(sv service.LocalityService) *LocalityDefault {
	return &LocalityDefault{sv: sv}
}

// LocalityDefault is a struct with methods that represent handlers for Locality
type LocalityDefault struct {
	// sv is the service that will be used by the handler
	sv service.LocalityService
}

// Create is a method that returns a handler for the route POST /localities
func (h *LocalityDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var inputLocality models.LocalityCreateRequest

		err := json.NewDecoder(r.Body).Decode(&inputLocality)

		if err != nil || !models.IsValidCreateRequestLocalities(inputLocality) {
			srvErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// process
		var locality = models.FromCreateRequestLocality(inputLocality)

		data, err := h.sv.Create(locality)
		if err != nil {
			srvErr := pkg.ServiceErrors[pkg.ErrConflict]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// response
		var result = models.ToLocalityDoc(data)

		writeResponse(w, http.StatusCreated, map[string]any{
			"data": result,
		}, nil)
	}
}

// SellersByLocality is a method that returns a handler for the route GET /localities/reportSellers/{id}
func (h *LocalityDefault) SellersByLocality() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		inputId := chi.URLParam(r, "id")

		id, err := strconv.Atoi(inputId)

		if err != nil {
			// Error de parsing del ID → BadRequest
			srvErr := pkg.ServiceErrors[pkg.ErrBadRequest]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// process
		data, err := h.sv.GetCantSellersByLocality(id)
		if err != nil {
			srvErr := pkg.ServiceErrors[pkg.ErrNotFound]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// response

		writeResponse(w, http.StatusCreated, map[string]any{
			"data": data,
		}, nil)
	}
}
