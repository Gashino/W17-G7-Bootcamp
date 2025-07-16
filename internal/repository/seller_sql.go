package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// SellerSql is a repository for managing seller data with SQL database
type SellerSql struct {
	db *sql.DB
}

// SQL queries constants
const (
	// SELECT queries
	querySelectAllSellers = `SELECT id, cid, company_name, address, telephone, locality_id FROM sellers`
	querySelectSellerById = `SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?`
	queryCheckSellerByCId = `SELECT id FROM sellers WHERE cid = ?`

	// INSERT queries
	queryInsertSeller = `INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)`

	// UPDATE queries
	queryUpdateSellerBase = `UPDATE sellers SET %s WHERE id = ?`

	// DELETE queries
	queryDeleteSeller = `DELETE FROM sellers WHERE id = ?`
)

// NewSellerSql creates a new seller repository with SQL database connection
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
			&seller.LocalityID,
		)
		if err != nil {
			return nil, err
		}
		result[seller.ID] = seller
	}
	return result, nil
}

// Create is a method that creates a Seller if it doesn't exist
func (r *SellerSql) Create(seller models.Seller) (models.Seller, error) {
	// Check if seller with same CId already exists
	var existingId int
	err := r.db.QueryRow(queryCheckSellerByCId, seller.CId).Scan(&existingId)

	if err == nil {
		// Seller with this CId already exists
		return models.Seller{}, pkg.ServiceErrors[pkg.ErrConflict]
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Some other error occurred
		return models.Seller{}, pkg.ServiceErrors[pkg.ErrNotFound]
	}

	// Insert new seller
	result, err := r.db.Exec(queryInsertSeller, seller.CId, seller.CompanyName, seller.Address, seller.Telephone, seller.LocalityID)
	if err != nil {
		return models.Seller{}, fmt.Errorf("failed to create seller: %w", err)
	}

	// Get the auto-generated ID
	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Seller{}, fmt.Errorf("failed to get last insert id: %w", err)
	}

	seller.ID = int(lastId)
	return seller, nil
}

// GetById is a method that returns a Seller by its ID if it exists
func (r *SellerSql) GetById(id int) (models.Seller, error) {
	var seller models.Seller
	err := r.db.QueryRow(querySelectSellerById, id).Scan(
		&seller.ID,
		&seller.CId,
		&seller.CompanyName,
		&seller.Address,
		&seller.Telephone,
		&seller.LocalityID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Seller{}, pkg.ServiceErrors[pkg.ErrNotFound]
		}
		return seller, err
	}
	return seller, nil
}

// UpdateFields is a method that updates a Seller's fields if it exists
func (r *SellerSql) UpdateFields(id int, data models.SellerCreateRequest) (models.Seller, error) {
	// First, check if seller exists
	_, err := r.GetById(id)
	if err != nil {
		return models.Seller{}, err
	}

	// Build dynamic update query
	setParts := []string{}
	args := []any{}

	if data.CId != nil {
		setParts = append(setParts, "cid = ?")
		args = append(args, *data.CId)
	}
	if data.CompanyName != nil {
		setParts = append(setParts, "company_name = ?")
		args = append(args, *data.CompanyName)
	}
	if data.Address != nil {
		setParts = append(setParts, "address = ?")
		args = append(args, *data.Address)
	}
	if data.Telephone != nil {
		setParts = append(setParts, "telephone = ?")
		args = append(args, *data.Telephone)
	}
	if data.LocalityID != nil {
		setParts = append(setParts, "locality_id = ?")
		args = append(args, *data.LocalityID)
	}

	if len(setParts) == 0 {
		// No fields to update, return current seller
		return r.GetById(id)
	}

	// Add id to args for WHERE clause
	args = append(args, id)

	// Build the update query properly
	updateQuery := fmt.Sprintf(queryUpdateSellerBase, strings.Join(setParts, ", "))

	// Execute update
	_, err = r.db.Exec(updateQuery, args...)
	if err != nil {
		return models.Seller{}, fmt.Errorf("failed to update seller: %w", err)
	}

	// Return updated seller
	return r.GetById(id)
}

// DeleteSeller is a method that deletes a Seller if it exists
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
