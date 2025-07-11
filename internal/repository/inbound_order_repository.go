package repository

import (
	"app/pkg/models"
)

type InboundOrderRepository interface {
	Create(inboundOrder models.InboundOrder) error
}
