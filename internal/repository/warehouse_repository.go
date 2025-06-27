package repository

import "app/pkg/models"

// interface that represents a warehouse repository
type WarehouseRepository interface {
	FindAll() (v map[int]models.Warehouse, err error)
}
