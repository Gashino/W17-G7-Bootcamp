package repository

import "app/pkg/models"

type ProductRecordRepository interface {
	Insert(record models.ProductRecord) (*models.ProductRecord, error)
}
