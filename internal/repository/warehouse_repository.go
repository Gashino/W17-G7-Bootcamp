package repository

import "app/pkg/models"

// interface that represents a warehouse repository
type WarehouseRepository interface {
	FindAll() (v map[int]models.Warehouse, err error)
	FindOne(id int) (v models.Warehouse, err error)
	Add(v models.Warehouse) (err error)
	FindAvailableID() (id int, err error)
	FindWarehouseByCode(code string) (v models.Warehouse, err error)
}
