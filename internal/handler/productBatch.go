package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bootcamp-go/web/response"
	"net/http"
)

// NewProductBatchDefault is a function that returns a new instance of ProductBatchDefault
func NewProductBatchDefault(sv service.ProductBatchService) *ProductBatchDefault {
	return &ProductBatchDefault{sv: sv}
}

// ProductBatchDefault is a struct with methods that represent handlers for sections
type ProductBatchDefault struct {
	// sv is the service that will be used by the handler
	sv service.ProductBatchService
}

func (h *ProductBatchDefault) CreateBatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...
		reqBody := models.ProductBatchPost{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			response.Error(w, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, pkg.ServiceErrors[pkg.ErrBadRequest].Error())
			return
		}

		batch := models.ProductBatch{
			ID: 0,
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        reqBody.BatchNumber,
				CurrentQuantity:    reqBody.CurrentQuantity,
				CurrentTemperature: reqBody.CurrentTemperature,
				DueDate:            reqBody.DueDate,
				InitialQuantity:    reqBody.InitialQuantity,
				ManufacturingDate:  reqBody.ManufacturingDate,
				ManufacturingHour:  reqBody.ManufacturingHour,
				MinumumTemperature: reqBody.MinumumTemperature,
				ProductId:          reqBody.ProductID,
				SectionId:          reqBody.SectionID,
			},
		}

		err = h.ValidatePostProductBatch(batch)
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
		value, err := h.sv.PostProductBatch(batch)
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
		data := models.ProductBatchDoc{
			ID:                 value.ID,
			BatchNumber:        value.BatchNumber,
			CurrentQuantity:    value.CurrentQuantity,
			CurrentTemperature: value.CurrentTemperature,
			DueDate:            value.DueDate,
			InitialQuantity:    value.InitialQuantity,
			ManufacturingDate:  value.ManufacturingDate,
			ManufacturingHour:  value.ManufacturingHour,
			MinumumTemperature: value.MinumumTemperature,
			ProductId:          value.ProductId,
			SectionId:          value.SectionId,
		}

		writeResponse(w, http.StatusCreated, map[string]any{
			"data": data,
		}, nil)
	}

}

func (h *ProductBatchDefault) ValidatePostProductBatch(batch models.ProductBatch) error {
	if batch.BatchNumber == 0 ||
		batch.CurrentQuantity == 0 ||
		batch.CurrentTemperature == 0.0 ||
		batch.DueDate == "" ||
		batch.InitialQuantity == 0 ||
		batch.ManufacturingDate == "" ||
		batch.ManufacturingHour == 0 ||
		batch.MinumumTemperature == 0.0 ||
		batch.ProductId == 0 ||
		batch.SectionId == 0 ||
		batch.CurrentTemperature < batch.MinumumTemperature {

		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		return svcErr
	}
	return nil
}
