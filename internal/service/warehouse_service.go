package service

import "app/pkg/models"

type WarehouseService interface {
	FindAll() (v map[int]models.Warehouse, err error)
}
