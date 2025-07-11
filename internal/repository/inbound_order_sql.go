package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type InboundOrderSQL struct {
	db *sql.DB
}

func NewInboundOrderSQL(db *sql.DB) *InboundOrderSQL {
	return &InboundOrderSQL{db: db}
}

func (r *InboundOrderSQL) Create(inboundOrder models.InboundOrder) error {
	_, err := r.db.Exec("INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES (?, ?, ?, ?, ?)", inboundOrder.OrderDate, inboundOrder.OrderNumber, inboundOrder.EmployeeID, inboundOrder.ProductBatchID, inboundOrder.WarehouseID)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1452:
				srvError := pkg.ServiceErrors[pkg.ErrNotFound]
				srvError.InternalError = fmt.Errorf(mysqlErr.Message)
				return srvError
			case 1062:
				srvError := pkg.ServiceErrors[pkg.ErrConflict]
				srvError.InternalError = fmt.Errorf(mysqlErr.Message)
				return srvError
			default:
				return pkg.ServiceErrors[pkg.ErrInternalServer]
			}
		}
		return pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	return nil
}
