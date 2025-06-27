package loader

import "app/pkg/models"

// SellerLoader is an interface that represents the loader for Sellers
type SellerLoader interface {
	// Load is a method that loads the Sellers
	Load() (s map[int]models.Seller, err error)
}
