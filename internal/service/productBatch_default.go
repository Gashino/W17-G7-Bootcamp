package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

// NewProductBatchDefault is a function that returns a new instance of ProductBatchDefault
func NewProductBatchDefault(rp repository.ProductBatchRepository) *ProductBatchDefault {
	return &ProductBatchDefault{rp: rp}
}

// ProductBatchDefault is a struct that represents the default service for productBatchs
type ProductBatchDefault struct {
	// rp is the repository that will be used by the service
	rp repository.ProductBatchRepository
}

func (sv *ProductBatchDefault) PostProductBatch(batch models.ProductBatch) (ProductBatch models.ProductBatch, err error) {
	return sv.rp.Create(section)
}
