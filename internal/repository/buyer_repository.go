package repository

import "app/pkg/models"

// RepositoryBuyer is an interface that defines the contract for buyer data operations
type RepositoryBuyer interface {
	// GetAll is a method that returns all buyers
	GetAll() (b map[int]models.Buyer, err error)
	// GetByID is a method that returns a buyer by its ID
	GetByID(id int) (b models.Buyer, err error)
	// Create is a method that creates a new buyer
	Create(buyer models.Buyer) (b models.Buyer, err error)
	// Update is a method that updates an existing buyer
	Update(buyer models.Buyer) (b models.Buyer, err error)
	// Delete is a method that deletes a buyer by its ID
	Delete(id int) (err error)
}
