package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

func NewProductRecordDefault(rp repository.ProductRecordRepository) *ProductRecordDefault {
	return &ProductRecordDefault{rp: rp}
}

type ProductRecordDefault struct {
	rp repository.ProductRecordRepository
}

func (p ProductRecordDefault) Create(productRecord models.ProductRecord) (*models.ProductRecord, error) {
	return p.rp.Insert(productRecord)
}
