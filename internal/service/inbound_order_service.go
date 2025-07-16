package service

import "app/pkg/models"

type InboundOrderService interface {
	Create(inboundOrder models.InboundOrder) (models.InboundOrder, error)
}
