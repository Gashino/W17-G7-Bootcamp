package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

// NewPurchaseOrderDefault is a function that returns a new instance of PurchaseOrderDefault
func NewPurchaseOrderDefault(rp repository.RepositoryPurchaseOrder) *PurchaseOrderDefault {
	return &PurchaseOrderDefault{
		rp: rp,
	}
}

// PurchaseOrderDefault is a struct that implements the ServicePurchaseOrder interface
type PurchaseOrderDefault struct {
	// rp is the repository for purchase order data operations
	rp repository.RepositoryPurchaseOrder
}

// Create is a method that creates a new purchase order
func (s *PurchaseOrderDefault) Create(purchaseOrder models.PurchaseOrder) (po models.PurchaseOrder, err error) {
	// Create the purchase order (this will assign an ID and check uniqueness)
	po, err = s.rp.Create(purchaseOrder)
	return
}
