package service

import (
	"app/pkg/models"
)

type ProductBatchService interface {
	PostProductBatch(batch models.ProductBatch) (ProductBatch models.ProductBatch, err error)
}
