package purchase_order

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockPurchaseOrderRepository struct {
	mock.Mock
}

func (m *MockPurchaseOrderRepository) Create(purchaseOrder models.PurchaseOrder) (models.PurchaseOrder, error) {
	args := m.Called(purchaseOrder)
	return args.Get(0).(models.PurchaseOrder), args.Error(1)
}
