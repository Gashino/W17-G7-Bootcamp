package repository

import (
	"app/pkg"
	"app/pkg/models"
)

// SellerMap is an in-memory repository for managing sellers
type SellerMap struct {
	db    map[int]models.Seller
	maxId int
}

// NewSellerMap creates a new seller repository with initial data
func NewSellerMap(db []models.SellerDoc) *SellerMap {
	var maxId int
	dbMap := make(map[int]models.Seller)

	for _, s := range db {
		if s.ID > maxId {
			maxId = s.ID
		}
		dbMap[s.ID] = models.FromSellerDoc(s)
	}

	return &SellerMap{
		db:    dbMap,
		maxId: maxId,
	}
}

// FindAll is a method that returns a map of all Sellers
func (r *SellerMap) FindAll() (v map[int]models.Seller, err error) {
	v = make(map[int]models.Seller)

	// copy db
	for key, value := range r.db {
		v[key] = value
	}

	return
}

// Create is a method that create a Seller if not exists
func (r *SellerMap) Create(seller models.Seller) (models.Seller, error) {

	largo := len(r.db) + 1
	_, find := r.db[seller.Id]
	if find {
		return models.Seller{}, pkg.ServiceErrors[pkg.ErrConflict]
	}
	seller.Id = largo
	r.db[largo] = seller
	return seller, nil
}

// GetById is a method that returns a Seller if exists
func (r *SellerMap) GetById(id int) (models.Seller, error) {
	result, find := r.db[id]
	if !find {
		return models.Seller{}, pkg.ServiceErrors[pkg.ErrNotFound]
	}
	return result, nil
}

func (r *SellerMap) UpdateFields(id int, data models.SellerCreateRequest) (models.Seller, error) {
	seller, find := r.db[id]
	if !find {
		return models.Seller{}, pkg.ServiceErrors[pkg.ErrNotFound]
	}
	if data.Address != nil {
		seller.Adress = *data.Address
	}
	if data.CId != nil {
		seller.CId = *data.CId
	}
	if data.Telephone != nil {
		seller.Telephone = *data.Telephone
	}
	if data.CompanyName != nil {
		seller.CompanyName = *data.CompanyName
	}
	r.db[id] = seller

	return seller, nil
}

func (r *SellerMap) DeleteSeller(id int) error {
	_, find := r.db[id]
	if !find {
		return pkg.ServiceErrors[pkg.ErrNotFound]
	}
	delete(r.db, id)
	return nil
}
