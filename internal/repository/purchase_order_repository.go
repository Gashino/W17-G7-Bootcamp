package repository

import "app/pkg/models"

// RepositoryPurchaseOrder is an interface that defines the contract for purchase order data operations
type RepositoryPurchaseOrder interface {
	// Create is a method that creates a new purchase order
	Create(purchaseOrder models.PurchaseOrder) (po models.PurchaseOrder, err error)
}
