package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bootcamp-go/web/response"
)

type InboundOrderHandler struct {
	service service.InboundOrderService
}

func NewInboundOrderHandler(service service.InboundOrderService) *InboundOrderHandler {
	return &InboundOrderHandler{service: service}
}

func (h *InboundOrderHandler) CreateInboundOrder(w http.ResponseWriter, r *http.Request) {
	var inboundOrder models.InboundOrder
	if err := json.NewDecoder(r.Body).Decode(&inboundOrder); err != nil {
		srvError := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		response.Error(w, srvError.ResponseCode, srvError.Error())
		return
	}

	// Validate the inbound order data
	if err := models.ValidateInboundOrder(inboundOrder, false); err != nil {
		srvError := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		srvError.InternalError = fmt.Errorf(err.Error())
		response.Error(w, http.StatusUnprocessableEntity, srvError.Error())
		return
	}

	order, err := h.service.Create(inboundOrder)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	data := map[string]any{
		"data": order.MapToDTO(),
	}

	writeResponse(w, http.StatusCreated, data, nil)
}
