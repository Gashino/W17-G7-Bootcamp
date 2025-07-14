package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// PurchaseOrderSQL is a struct that implements the RepositoryPurchaseOrder interface using SQL
type PurchaseOrderSQL struct {
	db *sql.DB
}

// NewPurchaseOrderSQL is a function that returns a new instance of PurchaseOrderSQL
func NewPurchaseOrderSQL(db *sql.DB) *PurchaseOrderSQL {
	return &PurchaseOrderSQL{
		db: db,
	}
}

// Create is a method that creates a new purchase order
func (r *PurchaseOrderSQL) Create(purchaseOrder models.PurchaseOrder) (po models.PurchaseOrder, err error) {
	result, err := r.db.Exec(
		"INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id) VALUES (?, ?, ?, ?, ?)",
		purchaseOrder.OrderNumber, purchaseOrder.OrderDate, purchaseOrder.TrackingCode, purchaseOrder.BuyerID, purchaseOrder.ProductRecordID,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				// Duplicate entry error - order_number already exists
				srvError := pkg.ServiceErrors[pkg.ErrBadRequest]
				srvError.InternalError = fmt.Errorf("order_number already exists")
				return models.PurchaseOrder{}, srvError
			case 1452:
				// Foreign key constraint error - buyer_id doesn't exist
				srvError := pkg.ServiceErrors[pkg.ErrBadRequest]
				srvError.InternalError = fmt.Errorf("buyer_id does not exist")
				return models.PurchaseOrder{}, srvError
			default:
				// Other MySQL error
				return models.PurchaseOrder{}, pkg.ServiceErrors[pkg.ErrInternalServer]
			}
		}
		// Non-MySQL error
		return models.PurchaseOrder{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return models.PurchaseOrder{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	purchaseOrder.ID = int(id)
	return purchaseOrder, nil
}
