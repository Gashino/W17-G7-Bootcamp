package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"errors"
	"fmt"
)

// LocalitySql is a repository for managing locality data with SQL database
type LocalitySql struct {
	db *sql.DB
}

// SQL queries constants
const (
	// SELECT queries
	querySelectLocalityById  = `SELECT id, locality_name, province_name, country_name FROM localities WHERE id = ?`
	queryCheckLocalityByName = `SELECT id FROM localities WHERE locality_name = ?`
	// INSERT queries
	queryInsertLocality = `INSERT INTO localities (locality_name, province_name, country_name) VALUES (?, ?, ?)`

	queryCantSellersByLocality = `SELECT l.id AS idLocalidad, l.locality_name AS nombreLocalidad, COUNT(s.id) AS cantidadSellers
	FROM localities l
	LEFT JOIN sellers s ON s.locality_id = l.id
	WHERE l.id = ?
	GROUP BY l.id, l.locality_name;`
)

// NewLocalitySql creates a new locality repository with initial data
func NewLocalitySql(db *sql.DB) *LocalitySql {
	return &LocalitySql{
		db: db,
	}
}

// Create is a method that creates a Locality if it doesn't exist
func (r *LocalitySql) Create(locality models.Locality) (models.Locality, error) {
	// Check if locality with same name already exists
	var existingId int
	err := r.db.QueryRow(queryCheckLocalityByName, locality.LocalityName).Scan(&existingId)

	if err == nil {
		// Locality with this name already exists
		return models.Locality{}, pkg.ServiceErrors[pkg.ErrConflict]
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Some other error occurred
		return models.Locality{}, pkg.ServiceErrors[pkg.ErrNotFound]
	}

	// Insert new locality
	result, err := r.db.Exec(queryInsertLocality, locality.LocalityName, locality.ProvinceName, locality.CountryName)
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

// GetById is a method that returns a Locality if it exists
func (r *LocalitySql) GetById(id int) (models.Locality, error) {
	var locality models.Locality
	err := r.db.QueryRow(querySelectLocalityById, id).Scan(
		&locality.ID,
		&locality.LocalityName,
		&locality.ProvinceName,
		&locality.CountryName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Locality{}, pkg.ServiceErrors[pkg.ErrNotFound]
		}
		return locality, err
	}
	return locality, nil
}

// GetCantSellersByLocality is a method that returns sellers count by locality
func (r *LocalitySql) GetCantSellersByLocality(id int) (models.LocalityBySellerResponse, error) {
	var locality models.LocalityBySellerResponse
	err := r.db.QueryRow(queryCantSellersByLocality, id).Scan(
		&locality.ID,
		&locality.LocalityName,
		&locality.SellerCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.LocalityBySellerResponse{}, pkg.ServiceErrors[pkg.ErrNotFound]
		}
		return locality, err
	}
	return locality, nil
}
