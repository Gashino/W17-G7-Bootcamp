package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"encoding/json"
	"net/http"

	"github.com/bootcamp-go/web/response"
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

// Create is a method that returns a handler for de route POST /Sellers
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
