package models

import "time"

type ProductRecord struct {
	LastUpdateDate *time.Time `json:"last_update_date"`
	PurchasePrice  *float64   `json:"purchase_price"`
	SalePrice      *float64   `json:"sale_price"`
	ProductId      *int       `json:"product_id"`
}
