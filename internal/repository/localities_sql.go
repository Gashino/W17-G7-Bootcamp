package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"errors"
	"fmt"
)

// LocalitySql is an in-memory repository for managing locality
type LocalitySql struct {
	db *sql.DB
}

// SQL queries constants
const (
	// SELECT queries
	querySelectLocalityById  = `SELECT id, locality_name, province_name, country_name FROM localities WHERE id = ?`
	queryCheckLocalityByName = `SELECT id FROM localities WHERE locality_name = ?`
	// INSERT queries
	queryInsertLocality = `INSERT INTO localities (cid, locality_name, province_name, country_name) VALUES (?, ?, ?, ?)`
)

// NewLocalitySql creates a new seller repository with initial data
func NewLocalitySql(db *sql.DB) *LocalitySql {
	return &LocalitySql{
		db: db,
	}
}

// Create is a method that create a Locality if not exists
func (r *LocalitySql) Create(locality models.Locality) (models.Locality, error) {
	// Check if seller with same CId already exists
	var existingId int
	err := r.db.QueryRow(queryCheckLocalityByName, locality.LocalityName).Scan(&existingId)

	if err == nil {
		// Seller with this CId already exists
		return models.Locality{}, pkg.ServiceErrors[pkg.ErrConflict]
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Some other error occurred
		return models.Locality{}, pkg.ServiceErrors[pkg.ErrNotFound]
	}

	// Insert new seller
	result, err := r.db.Exec(queryInsertLocality, locality.CountryName, locality.LocalityName, locality.ProvinceName)
	if err != nil {
		return models.Locality{}, fmt.Errorf("failed to create Locality: %w", err)
	}

	// Get the auto-generated ID
	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Locality{}, fmt.Errorf("failed to get last insert id: %w", err)
	}

	locality.ID = int(lastId)
	return locality, nil
}

// GetById is a method that returns a Seller if exists
func (r *LocalitySql) GetById(id int) (models.Locality, error) {
	var locality models.Locality
	err := r.db.QueryRow(querySelectLocalityById, id).Scan(
		&locality.ID,
		&locality.CountryName,
		&locality.LocalityName,
		&locality.ProvinceName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Locality{}, pkg.ServiceErrors[pkg.ErrNotFound]
		}
		return locality, err
	}
	return locality, nil
}
