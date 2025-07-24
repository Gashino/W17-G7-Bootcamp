package inbound_order

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockInboundOrderRepository struct {
	mock.Mock
}

// Create implements repository.InboundOrderRepository.
func (m *MockInboundOrderRepository) Create(inboundOrder models.InboundOrder) (models.InboundOrder, error) {
	args := m.Called(inboundOrder)
	return args.Get(0).(models.InboundOrder), args.Error(1)
}
