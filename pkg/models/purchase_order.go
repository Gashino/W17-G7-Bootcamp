package models

// PurchaseOrderDoc is a struct that represents a purchase order document in JSON format
type PurchaseOrderDoc struct {
	ID              int    `json:"id"`
	OrderNumber     string `json:"order_number"`
	OrderDate       string `json:"order_date"`
	TrackingCode    string `json:"tracking_code"`
	BuyerID         int    `json:"buyer_id"`
	ProductRecordID int    `json:"product_record_id"`
}

// PurchaseOrder is a struct that represents a purchase order
type PurchaseOrder struct {
	// ID is the unique identifier of the purchase order
	ID int `json:"id"`

	// PurchaseOrderAttributes contains the attributes of a purchase order
	PurchaseOrderAttributes
}

// PurchaseOrderAttributes is a struct that contains the attributes of a purchase order
type PurchaseOrderAttributes struct {
	OrderNumber     string `json:"order_number"`
	OrderDate       string `json:"order_date"`
	TrackingCode    string `json:"tracking_code"`
	BuyerID         int    `json:"buyer_id"`
	ProductRecordID int    `json:"product_record_id"`
}
