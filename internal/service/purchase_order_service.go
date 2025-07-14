package service

import "app/pkg/models"

// ServicePurchaseOrder is an interface that defines the contract for purchase order business logic operations
type ServicePurchaseOrder interface {
	// Create is a method that creates a new purchase order
	Create(purchaseOrder models.PurchaseOrder) (po models.PurchaseOrder, err error)
}
