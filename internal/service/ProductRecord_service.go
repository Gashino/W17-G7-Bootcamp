package service

import "app/pkg/models"

type IProductRecordService interface {
	Create(productRecord models.ProductRecord) (*models.ProductRecord, error)
}
