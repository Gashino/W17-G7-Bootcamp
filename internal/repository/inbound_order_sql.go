package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const (
	insertInboundOrder = "INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES (?, ?, ?, ?, ?)"
)

type InboundOrderSQL struct {
	db *sql.DB
}

func NewInboundOrderSQL(db *sql.DB) *InboundOrderSQL {
	return &InboundOrderSQL{db: db}
}

func (r *InboundOrderSQL) Create(inboundOrder models.InboundOrder) (models.InboundOrder, error) {
	_, err := r.db.Exec(insertInboundOrder, inboundOrder.OrderDate, inboundOrder.OrderNumber, inboundOrder.EmployeeID, inboundOrder.ProductBatchID, inboundOrder.WarehouseID)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1452:
				srvError := pkg.ServiceErrors[pkg.ErrNotFound]
				srvError.InternalError = fmt.Errorf(mysqlErr.Message)
				return models.InboundOrder{}, srvError
			case 1062:
				srvError := pkg.ServiceErrors[pkg.ErrConflict]
				srvError.InternalError = fmt.Errorf(mysqlErr.Message)
				return models.InboundOrder{}, srvError
			default:
				return models.InboundOrder{}, pkg.ServiceErrors[pkg.ErrInternalServer]
			}
		}
		return models.InboundOrder{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	return inboundOrder, nil
}
