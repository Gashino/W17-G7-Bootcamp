package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// NewSectionDefault is a function that returns a new instance of SectionDefault
func NewSectionDefault(sv service.SectionService) *SectionDefault {
	return &SectionDefault{sv: sv}
}

// SectionDefault is a struct with methods that represent handlers for sections
type SectionDefault struct {
	// sv is the service that will be used by the handler
	sv service.SectionService
}

// GetAll is a method that returns a handler for the route GET /sections
func (h *SectionDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		// process
		// - get all sections
		v, err := h.sv.GetAll()
		if err != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(err, &svcErr) {
				response.Error(w, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(w, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		// response
		data := []models.SectionDoc{}
		for _, value := range v {
			sec := models.SectionDoc{
				ID:                 value.ID,
				SectionNumber:      value.SectionNumber,
				CurrentTemperature: value.CurrentTemperature,
				MinimumTemperature: value.MinimumTemperature,
				CurrentCapacity:    value.CurrentCapacity,
				MinimumCapacity:    value.MinimumCapacity,
				MaximumCapacity:    value.MaximumCapacity,
				WarehouseID:        value.WarehouseID,
				ProductTypeID:      value.ProductTypeID,
			}
			data = append(data, sec)
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *SectionDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...
		idString := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			response.Error(w, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, pkg.ServiceErrors[pkg.ErrBadRequest].Error())
			return
		}

		// process
		// - get all sections
		value, err := h.sv.GetByID(id)
		if err != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(err, &svcErr) {
				response.Error(w, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(w, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		// response
		data := models.SectionDoc{
			ID:                 value.ID,
			SectionNumber:      value.SectionNumber,
			CurrentTemperature: value.CurrentTemperature,
			MinimumTemperature: value.MinimumTemperature,
			CurrentCapacity:    value.CurrentCapacity,
			MinimumCapacity:    value.MinimumCapacity,
			MaximumCapacity:    value.MaximumCapacity,
			WarehouseID:        value.WarehouseID,
			ProductTypeID:      value.ProductTypeID,
		}

		writeResponse(w, http.StatusOK, map[string]any{
			"data": data,
		}, nil)
	}

}

func (h *SectionDefault) PostSection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...
		reqBody := models.SectionPostDoc{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			response.Error(w, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, pkg.ServiceErrors[pkg.ErrBadRequest].Error())
			return
		}

		section := models.Section{
			ID: 0,
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      reqBody.SectionNumber,
				CurrentTemperature: reqBody.CurrentTemperature,
				MinimumTemperature: reqBody.MinimumTemperature,
				CurrentCapacity:    reqBody.CurrentCapacity,
				MinimumCapacity:    reqBody.MinimumCapacity,
				MaximumCapacity:    reqBody.MaximumCapacity,
				WarehouseID:        reqBody.WarehouseID,
				ProductTypeID:      reqBody.ProductTypeID,
				ProductBatches:     nil,
			},
		}

		err = h.ValidatePostSection(section)
		if err != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(err, &svcErr) {
				response.Error(w, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(w, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		// process
		// - get all sections
		value, err := h.sv.Create(section)
		if err != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(err, &svcErr) {
				response.Error(w, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(w, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		// response
		data := models.SectionDoc{
			ID:                 value.ID,
			SectionNumber:      value.SectionNumber,
			CurrentTemperature: value.CurrentTemperature,
			MinimumTemperature: value.MinimumTemperature,
			CurrentCapacity:    value.CurrentCapacity,
			MinimumCapacity:    value.MinimumCapacity,
			MaximumCapacity:    value.MaximumCapacity,
			WarehouseID:        value.WarehouseID,
			ProductTypeID:      value.ProductTypeID,
		}

		writeResponse(w, http.StatusCreated, map[string]any{
			"data": data,
		}, nil)
	}

}

func (h *SectionDefault) ValidatePostSection(section models.Section) error {
	if section.SectionNumber == 0 ||
		section.CurrentTemperature == 0.0 ||
		section.MinimumTemperature == 0.0 ||
		section.CurrentCapacity == 0 ||
		section.MinimumCapacity == 0 ||
		section.MaximumCapacity == 0 ||
		section.WarehouseID == 0 ||
		section.ProductTypeID == 0 ||
		section.CurrentCapacity > section.MaximumCapacity ||
		section.MinimumCapacity > section.MaximumCapacity ||
		section.CurrentTemperature < section.MinimumTemperature {

		svcErr := pkg.ServiceErrors[pkg.ErrBadRequest]
		svcErr.InternalError = fmt.Errorf("section is not valid")
		return svcErr
	}
	return nil
}

func (h *SectionDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		idString := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			response.Error(w, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, pkg.ServiceErrors[pkg.ErrBadRequest].Error())
			return
		}

		reqBody := models.SectionPostDoc{}
		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			response.Error(w, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, pkg.ServiceErrors[pkg.ErrBadRequest].Error())
			return
		}

		section := models.Section{
			ID: 0,
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      reqBody.SectionNumber,
				CurrentTemperature: reqBody.CurrentTemperature,
				MinimumTemperature: reqBody.MinimumTemperature,
				CurrentCapacity:    reqBody.CurrentCapacity,
				MinimumCapacity:    reqBody.MinimumCapacity,
				MaximumCapacity:    reqBody.MaximumCapacity,
				WarehouseID:        reqBody.WarehouseID,
				ProductTypeID:      reqBody.ProductTypeID,
				ProductBatches:     nil,
			},
		}

		err = h.ValidatePatchSection(section)
		if err != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(err, &svcErr) {
				response.Error(w, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(w, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		// process
		// - get all sections
		value, err := h.sv.Update(id, section)
		if err != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(err, &svcErr) {
				response.Error(w, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(w, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		// response
		data := models.SectionDoc{
			ID:                 value.ID,
			SectionNumber:      value.SectionNumber,
			CurrentTemperature: value.CurrentTemperature,
			MinimumTemperature: value.MinimumTemperature,
			CurrentCapacity:    value.CurrentCapacity,
			MinimumCapacity:    value.MinimumCapacity,
			MaximumCapacity:    value.MaximumCapacity,
			WarehouseID:        value.WarehouseID,
			ProductTypeID:      value.ProductTypeID,
		}

		writeResponse(w, http.StatusCreated, map[string]any{
			"data": data,
		}, nil)
	}
}

func (h *SectionDefault) ValidatePatchSection(section models.Section) error {
	// Verificar que al menos un campo venga en el request
	if section.SectionNumber == 0 &&
		section.CurrentTemperature == 0.0 &&
		section.MinimumTemperature == 0.0 &&
		section.CurrentCapacity == 0 &&
		section.MinimumCapacity == 0 &&
		section.MaximumCapacity == 0 &&
		section.WarehouseID == 0 &&
		section.ProductTypeID == 0 {

		svcErr := pkg.ServiceErrors[pkg.ErrBadRequest]
		svcErr.InternalError = fmt.Errorf("at least one field must be provided")
		return svcErr
	}
	return nil
}

func (h *SectionDefault) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		idString := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			response.Error(w, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, pkg.ServiceErrors[pkg.ErrBadRequest].Error())
			return
		}

		// process
		// - get all sections
		err = h.sv.Delete(id)
		if err != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(err, &svcErr) {
				response.Error(w, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(w, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		// response
		data := models.SectionDoc{}

		response.JSON(w, http.StatusNoContent, map[string]any{
			"message": "Deleted successfully",
			"data":    data,
		})
	}
}

func (h *SectionDefault) ReportProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...
		idString := r.URL.Query().Get("id")
		var idPtr *int

		if idString != "" {
			idInt, err := strconv.Atoi(idString)
			if err != nil {
				response.Error(w, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, pkg.ServiceErrors[pkg.ErrBadRequest].Error())
				return
			}
			idPtr = &idInt
		}

		// process
		// - get all sections
		v, err := h.sv.ReportProductsBySection(idPtr)
		if err != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(err, &svcErr) {
				response.Error(w, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(w, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		// response
		data := v
		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}

}
