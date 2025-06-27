package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

// NewBuyerDefault is a function that returns a new instance of BuyerDefault
func NewBuyerDefault(rp repository.RepositoryBuyer) *BuyerDefault {
	return &BuyerDefault{
		rp: rp,
	}
}

// BuyerDefault is a struct that implements the ServiceBuyer interface
type BuyerDefault struct {
	// rp is the repository for buyer data operations
	rp repository.RepositoryBuyer
}

// GetAll is a method that returns all buyers
func (s *BuyerDefault) GetAll() (b map[int]models.Buyer, err error) {
	b, err = s.rp.GetAll()
	return
}

// GetByID is a method that returns a buyer by its ID
func (s *BuyerDefault) GetByID(id int) (b models.Buyer, err error) {
	b, err = s.rp.GetByID(id)
	return
}

// Create is a method that creates a new buyer
func (s *BuyerDefault) Create(buyer models.Buyer) (b models.Buyer, err error) {
	// Create the buyer (this will assign an ID and check uniqueness)
	b, err = s.rp.Create(buyer)
	return
}

// Update is a method that updates an existing buyer
func (s *BuyerDefault) Update(id int, buyer models.Buyer) (b models.Buyer, err error) {
	// Set the ID from the path parameter
	buyer.Id = id

	// Update the buyer (this will handle partial updates and uniqueness validation)
	b, err = s.rp.Update(buyer)
	return
}

// Delete is a method that deletes a buyer by its ID
func (s *BuyerDefault) Delete(id int) (err error) {
	err = s.rp.Delete(id)
	return
}
