package repository

import (
	"app/pkg"
	"app/pkg/models"
	"sync"
)

// NewBuyerMap is a function that returns a new instance of BuyerMap
func NewBuyerMap(db map[int]models.Buyer) *BuyerMap {
	return &BuyerMap{
		db:     db,
		nextID: getNextID(db),
		mutex:  &sync.RWMutex{},
	}
}

// BuyerMap is a struct that implements the RepositoryBuyer interface using a map
type BuyerMap struct {
	// db is the map that stores the buyers
	db map[int]models.Buyer
	// nextID is the next available ID for new buyers
	nextID int
	// mutex is used for thread safety
	mutex *sync.RWMutex
}

// GetAll is a method that returns all buyers
func (r *BuyerMap) GetAll() (b map[int]models.Buyer, err error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	b = make(map[int]models.Buyer)
	for id, buyer := range r.db {
		b[id] = buyer
	}
	return
}

// GetByID is a method that returns a buyer by its ID
func (r *BuyerMap) GetByID(id int) (b models.Buyer, err error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	buyer, exists := r.db[id]
	if !exists {
		err = pkg.ServiceErrors[pkg.ErrNotFound]
		return
	}
	b = buyer
	return
}

// Create is a method that creates a new buyer
func (r *BuyerMap) Create(buyer models.Buyer) (b models.Buyer, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if card_number_id already exists
	for _, existingBuyer := range r.db {
		if existingBuyer.CardNumberID == buyer.CardNumberID {
			err = pkg.ServiceErrors[pkg.ErrBadRequest]
			return
		}
	}

	// Assign the next available ID
	buyer.Id = r.nextID
	r.nextID++

	// Store the buyer
	r.db[buyer.Id] = buyer

	// Return the created buyer
	b = buyer
	return
}

// Update is a method that updates an existing buyer
func (r *BuyerMap) Update(buyer models.Buyer) (b models.Buyer, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if buyer exists
	existingBuyer, exists := r.db[buyer.Id]
	if !exists {
		err = pkg.ServiceErrors[pkg.ErrNotFound]
		return
	}

	// If card_number_id is being changed, check for uniqueness
	if buyer.CardNumberID != "" && buyer.CardNumberID != existingBuyer.CardNumberID {
		for id, otherBuyer := range r.db {
			if id != buyer.Id && otherBuyer.CardNumberID == buyer.CardNumberID {
				err = pkg.ServiceErrors[pkg.ErrBadRequest]
				return
			}
		}
	}

	// Apply partial updates - only update fields that are not empty
	if buyer.CardNumberID != "" {
		existingBuyer.CardNumberID = buyer.CardNumberID
	}
	if buyer.FirstName != "" {
		existingBuyer.FirstName = buyer.FirstName
	}
	if buyer.LastName != "" {
		existingBuyer.LastName = buyer.LastName
	}

	// Store the updated buyer
	r.db[buyer.Id] = existingBuyer

	// Return the updated buyer
	b = existingBuyer
	return
}

// Delete is a method that deletes a buyer by its ID
func (r *BuyerMap) Delete(id int) (err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.db[id]
	if !exists {
		err = pkg.ServiceErrors[pkg.ErrNotFound]
		return
	}

	delete(r.db, id)
	return
}

// getNextID is a helper function that returns the next available ID
func getNextID(db map[int]models.Buyer) int {
	maxID := 0
	for id := range db {
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1
}
