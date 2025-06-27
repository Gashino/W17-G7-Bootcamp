package loader

import "app/pkg/models"

// LoaderBuyer is an interface that defines the contract for loading buyers
type LoaderBuyer interface {
	// Load is a method that loads the buyers
	Load() (b map[int]models.Buyer, err error)
}
