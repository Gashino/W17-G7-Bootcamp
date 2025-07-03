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

// NewSellerDefault is a function that returns a new instance of SellerDefault
func NewSellerDefault(sv service.SellerService) *SellerDefault {
	return &SellerDefault{sv: sv}
}

// SellerDefault is a struct with methods that represent handlers for Sellers
type SellerDefault struct {
	// sv is the service that will be used by the handler
	sv service.SellerService
}

// GetAll is a method that returns a handler for the route GET /Sellers
func (h *SellerDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		// process
		// - get all Sellers
		v, err := h.sv.FindAll()
		if err != nil {
			srvErr := pkg.ServiceErrors[pkg.ErrInternalServer]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// response
		data := make(map[int]models.SellerDoc)
		for key, value := range v {
			data[key] = models.ToSellerDoc(value)
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// GetById is a method that returns a handler for the route GET /Sellers/Id
func (h *SellerDefault) GetById() http.HandlerFunc {
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
		data, err := h.sv.GetById(id)
		if err != nil {
			srvErr := pkg.ServiceErrors[pkg.ErrNotFound]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// response
		var result = models.ToSellerDoc(data)

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    result,
		})
	}
}

// Create is a method that returns a handler for de route POST /Sellers
func (h *SellerDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var inputSeller models.SellerCreateRequest

		err := json.NewDecoder(r.Body).Decode(&inputSeller)

		if err != nil || !models.IsValidCreateRequest(inputSeller) {
			srvErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// process
		var seller = models.FromCreateRequest(inputSeller)

		data, err := h.sv.Create(seller)
		if err != nil {
			srvErr := pkg.ServiceErrors[pkg.ErrConflict]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// response
		var result = models.ToSellerDoc(data)

		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    result,
		})
	}
}

// Update is a method that returns a handler for de route PATCH /Sellers/Id
func (h *SellerDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var inputSeller models.SellerCreateRequest
		inputId := chi.URLParam(r, "id")

		id, errId := strconv.Atoi(inputId)

		err := json.NewDecoder(r.Body).Decode(&inputSeller)

		if err != nil || errId != nil {
			srvErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// process

		data, err := h.sv.UpdateFields(id, inputSeller)
		if err != nil {
			srvErr := pkg.ServiceErrors[pkg.ErrNotFound]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// response
		var result = models.ToSellerDoc(data)

		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    result,
		})
	}
}

// Delete is a method that returns a handler for de route DELETE /Sellers/Id
func (h *SellerDefault) Delete() http.HandlerFunc {
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
		err = h.sv.DeleteSeller(id)
		if err != nil {
			srvErr := pkg.ServiceErrors[pkg.ErrNotFound]
			response.Error(w, srvErr.ResponseCode, srvErr.Error())
			return
		}

		// response

		response.JSON(w, http.StatusNoContent, map[string]any{
			"message": "success",
			"data":    nil,
		})
	}
}
