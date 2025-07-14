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
	return s.repository.Create(inboundOrder)
}
