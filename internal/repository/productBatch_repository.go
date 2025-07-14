package repository

import "app/pkg/models"

type ProductBatchRepository interface {
	InsertProductBatch(pb *models.ProductBatch) (productBatch *models.ProductBatch, err error)
}
