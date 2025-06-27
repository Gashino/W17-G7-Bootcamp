package service

import "app/pkg/models"

type WarehouseService interface {
	FindAll() (v map[int]models.Warehouse, err error)
	FindOne(id int) (v models.Warehouse, err error)
	Add(v models.WarehouseDoc) (err error)
}
