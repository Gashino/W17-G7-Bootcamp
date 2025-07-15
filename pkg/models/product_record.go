package models

import "time"

type ProductRecord struct {
	ID             int        `json:"id"`
	LastUpdateDate *time.Time `json:"last_update_date"`
	PurchasePrice  *float64   `json:"purchase_price"`
	SalePrice      *float64   `json:"sale_price"`
	ProductId      *int       `json:"product_id"`
}

type ProductRecordResponse struct {
	ProductId    int    `json:"product_id"`
	Description  string `json:"description"`
	RecordsCount int    `json:"records_count"`
}

func (data *ProductRecord) Validate() bool {
	if data.LastUpdateDate == nil {
		return false
	}
	if data.PurchasePrice == nil {
		return false
	}
	if data.SalePrice == nil {
		return false
	}
	if data.ProductId == nil {
		return false
	}
	return true
}
