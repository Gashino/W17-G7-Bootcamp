package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

type InboundOrderServiceDefault struct {
	repository repository.InboundOrderRepository
}

func NewInboundOrderServiceDefault(repository repository.InboundOrderRepository) InboundOrderService {
	return &InboundOrderServiceDefault{repository: repository}
}

func (s *InboundOrderServiceDefault) Create(inboundOrder models.InboundOrder) (models.InboundOrder, error) {
	// Validate the inbound order first
	if err := models.ValidateInboundOrder(inboundOrder, false); err != nil {
		return models.InboundOrder{}, err
	}

	return s.repository.Create(inboundOrder)
}
