package models

import "time"

type InboundOrder struct {
	ID             int       `json:"id"`
	OrderDate      time.Time `json:"order_date"`
	OrderNumber    string    `json:"order_number"`
	EmployeeID     int       `json:"employee_id"`
	ProductBatchID int       `json:"product_batch_id"`
	WarehouseID    int       `json:"warehouse_id"`
}

func (i InboundOrder) MapToDTO() InboundOrderDTO {
	return InboundOrderDTO{
		OrderDate:      i.OrderDate,
		OrderNumber:    i.OrderNumber,
		EmployeeID:     i.EmployeeID,
		ProductBatchID: i.ProductBatchID,
		WarehouseID:    i.WarehouseID,
	}
}

type InboundOrderDTO struct {
	OrderDate      time.Time `json:"order_date"`
	OrderNumber    string    `json:"order_number"`
	EmployeeID     int       `json:"employee_id"`
	ProductBatchID int       `json:"product_batch_id"`
	WarehouseID    int       `json:"warehouse_id"`
}
