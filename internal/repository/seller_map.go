package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"errors"
)

// SellerSql is an in-memory repository for managing sellers
type SellerSql struct {
	db *sql.DB
}

// SQL queries constants
const (
	// SELECT queries
	querySelectAllSellers = `SELECT id, cid, company_name, address, telephone FROM sellers`
	querySelectSellerById = `SELECT id, cid, company_name, address, telephone FROM sellers WHERE id = ?`
	// DELETE queries
	queryDeleteSeller = `DELETE FROM sellers WHERE id = ?`
)

// NewSellerSql creates a new seller repository with initial data
func NewSellerSql(db *sql.DB) *SellerSql {
	return &SellerSql{
		db: db,
	}
}

// FindAll is a method that returns a map of all Sellers
func (r *SellerSql) FindAll() (v map[int]models.Seller, err error) {
	rows, err := r.db.Query(querySelectAllSellers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[int]models.Seller)
	for rows.Next() {
		var seller models.Seller
		err := rows.Scan(
			&seller.ID,
			&seller.CId,
			&seller.CompanyName,
			&seller.Address,
			&seller.Telephone,
		)
		if err != nil {
			return nil, err
		}
		result[seller.ID] = seller
	}
	return result, nil
}

// Create is a method that create a Seller if not exists
func (r *SellerSql) Create(seller models.Seller) (models.Seller, error) {
	return seller, nil
}

// GetById is a method that returns a Seller if exists
func (r *SellerSql) GetById(id int) (models.Seller, error) {
	var seller models.Seller
	err := r.db.QueryRow(querySelectSellerById, id).Scan(
		&seller.ID,
		&seller.CId,
		&seller.CompanyName,
		&seller.Address,
		&seller.Telephone,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Seller{}, pkg.ServiceErrors[pkg.ErrNotFound]
		}
		return seller, err
	}
	return seller, nil
}

// UpdateFields is a method that modify a Seller if exists
func (r *SellerSql) UpdateFields(id int, data models.SellerCreateRequest) (models.Seller, error) {
	var seller models.Seller
	return seller, nil
}

// DeleteSeller is a method that delete a Seller if exists
func (r *SellerSql) DeleteSeller(id int) error {
	_, err := r.GetById(id)
	if err != nil {
		return err
	}
	result, err := r.db.Exec(queryDeleteSeller, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ServiceErrors[pkg.ErrNotFound]
	}

	return nil
}
