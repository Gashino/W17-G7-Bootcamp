package models

import (
	"fmt"
	"strings"
	"time"
)

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

func ValidateInboundOrder(inboundOrder InboundOrder, validateID bool) error {
	// Validate OrderDate is not zero
	if inboundOrder.OrderDate.IsZero() {
		return fmt.Errorf("order_date is required")
	}

	// Validate OrderDate is not in the future
	if inboundOrder.OrderDate.After(time.Now()) {
		return fmt.Errorf("order_date cannot be in the future")
	}

	// Validate OrderNumber is not empty
	if len(strings.TrimSpace(inboundOrder.OrderNumber)) == 0 {
		return fmt.Errorf("order_number is required")
	}

	// Validate EmployeeID is valid
	if inboundOrder.EmployeeID <= 0 {
		return fmt.Errorf("employee_id must be a valid positive number")
	}

	// Validate ProductBatchID is valid
	if inboundOrder.ProductBatchID <= 0 {
		return fmt.Errorf("product_batch_id must be a valid positive number")
	}

	// Validate WarehouseID is valid
	if inboundOrder.WarehouseID <= 0 {
		return fmt.Errorf("warehouse_id must be a valid positive number")
	}

	// Validate ID if required (for updates)
	if validateID {
		if inboundOrder.ID <= 0 {
			return fmt.Errorf("id is required and must be a valid positive number")
		}
	}

	return nil
}

// Legacy function for backward compatibility
func (i InboundOrder) Validate() error {
	return ValidateInboundOrder(i, false)
}
