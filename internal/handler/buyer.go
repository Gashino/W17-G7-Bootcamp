package handler

import "net/http"

// HandlerBuyer is an interface that defines the contract for buyer HTTP request handling
type HandlerBuyer interface {
	// GetAll is a method that handles GET requests to retrieve all buyers
	GetAll() http.HandlerFunc
	// GetByID is a method that handles GET requests to retrieve a buyer by ID
	GetByID() http.HandlerFunc
	// Create is a method that handles POST requests to create a new buyer
	Create() http.HandlerFunc
	// Update is a method that handles PATCH requests to update an existing buyer
	Update() http.HandlerFunc
	// Delete is a method that handles DELETE requests to delete a buyer
	Delete() http.HandlerFunc
}
