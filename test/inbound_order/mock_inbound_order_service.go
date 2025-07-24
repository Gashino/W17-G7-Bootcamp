package inbound_order

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockInboundOrderService struct {
	mock.Mock
}

// Create implements service.InboundOrderService.
func (m *MockInboundOrderService) Create(inboundOrder models.InboundOrder) (models.InboundOrder, error) {
	args := m.Called(inboundOrder)
	return args.Get(0).(models.InboundOrder), args.Error(1)
}
