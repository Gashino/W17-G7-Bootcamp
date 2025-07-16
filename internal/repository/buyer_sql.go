package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// BuyerSQL is a struct that implements the RepositoryBuyer interface using SQL
type BuyerSQL struct {
	db *sql.DB
}

// NewBuyerSQL is a function that returns a new instance of BuyerSQL
func NewBuyerSQL(db *sql.DB) *BuyerSQL {
	return &BuyerSQL{
		db: db,
	}
}

// GetAll is a method that returns all buyers
func (r *BuyerSQL) GetAll() (b map[int]models.Buyer, err error) {
	rows, err := r.db.Query("SELECT id, card_number_id, first_name, last_name FROM buyers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	b = make(map[int]models.Buyer)
	for rows.Next() {
		var buyer models.Buyer
		err := rows.Scan(&buyer.ID, &buyer.CardNumberID, &buyer.FirstName, &buyer.LastName)
		if err != nil {
			return nil, err
		}
		b[buyer.ID] = buyer
	}

	return b, nil
}

// GetByID is a method that returns a buyer by its ID
func (r *BuyerSQL) GetByID(id int) (b models.Buyer, err error) {
	row := r.db.QueryRow("SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?", id)

	err = row.Scan(&b.ID, &b.CardNumberID, &b.FirstName, &b.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Buyer{}, pkg.ServiceErrors[pkg.ErrNotFound]
		}
		return models.Buyer{}, err
	}

	return b, nil
}

// Create is a method that creates a new buyer
func (r *BuyerSQL) Create(buyer models.Buyer) (b models.Buyer, err error) {
	result, err := r.db.Exec(
		"INSERT INTO buyers (card_number_id, first_name, last_name) VALUES (?, ?, ?)",
		buyer.CardNumberID, buyer.FirstName, buyer.LastName,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				// Duplicate entry error - card_number_id already exists
				srvError := pkg.ServiceErrors[pkg.ErrBadRequest]
				srvError.InternalError = fmt.Errorf("card_number_id already exists")
				return models.Buyer{}, srvError
			default:
				// Other MySQL error
				return models.Buyer{}, pkg.ServiceErrors[pkg.ErrInternalServer]
			}
		}
		// Non-MySQL error
		return models.Buyer{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return models.Buyer{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	buyer.ID = int(id)
	return buyer, nil
}

// Update is a method that updates an existing buyer
func (r *BuyerSQL) Update(buyer models.Buyer) (b models.Buyer, err error) {
	// First check if the buyer exists
	_, err = r.GetByID(buyer.ID)
	if err != nil {
		return models.Buyer{}, err
	}

	// Build dynamic update query based on non-empty fields
	query := "UPDATE buyers SET "
	args := []interface{}{}
	updates := []string{}

	if buyer.CardNumberID != "" {
		updates = append(updates, "card_number_id = ?")
		args = append(args, buyer.CardNumberID)
	}
	if buyer.FirstName != "" {
		updates = append(updates, "first_name = ?")
		args = append(args, buyer.FirstName)
	}
	if buyer.LastName != "" {
		updates = append(updates, "last_name = ?")
		args = append(args, buyer.LastName)
	}

	// If no fields to update, return the current buyer
	if len(updates) == 0 {
		return r.GetByID(buyer.ID)
	}

	query += updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = ?"
	args = append(args, buyer.ID)

	_, err = r.db.Exec(query, args...)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				// Duplicate entry error - card_number_id already exists
				srvError := pkg.ServiceErrors[pkg.ErrBadRequest]
				srvError.InternalError = fmt.Errorf("card_number_id already exists")
				return models.Buyer{}, srvError
			default:
				// Other MySQL error
				return models.Buyer{}, pkg.ServiceErrors[pkg.ErrInternalServer]
			}
		}
		// Non-MySQL error
		return models.Buyer{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	// Return the updated buyer
	return r.GetByID(buyer.ID)
}

// Delete is a method that deletes a buyer by its ID
func (r *BuyerSQL) Delete(id int) (err error) {
	result, err := r.db.Exec("DELETE FROM buyers WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pkg.ServiceError{
			ResponseCode:  pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode,
			InternalError: pkg.ServiceErrors[pkg.ErrNotFound].InternalError,
		}
	}

	return nil
}

// GetPurchaseOrdersReport is a method that returns a purchase order report for all buyers or a specific buyer
func (r *BuyerSQL) GetPurchaseOrdersReport(buyerID *int) (reports []models.BuyerPurchaseOrderReport, err error) {
	var query string
	var args []interface{}

	if buyerID != nil {
		// Query for specific buyer
		query = `
			SELECT 
				b.id, 
				b.card_number_id, 
				b.first_name, 
				b.last_name, 
				COUNT(po.id) as purchase_orders_count
			FROM buyers b
			LEFT JOIN purchase_orders po ON b.id = po.buyer_id
			WHERE b.id = ?
			GROUP BY b.id, b.card_number_id, b.first_name, b.last_name
		`
		args = append(args, *buyerID)
	} else {
		// Query for all buyers
		query = `
			SELECT 
				b.id, 
				b.card_number_id, 
				b.first_name, 
				b.last_name, 
				COUNT(po.id) as purchase_orders_count
			FROM buyers b
			LEFT JOIN purchase_orders po ON b.id = po.buyer_id
			GROUP BY b.id, b.card_number_id, b.first_name, b.last_name
			ORDER BY b.id
		`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports = make([]models.BuyerPurchaseOrderReport, 0)
	for rows.Next() {
		var report models.BuyerPurchaseOrderReport
		err := rows.Scan(
			&report.ID,
			&report.CardNumberID,
			&report.FirstName,
			&report.LastName,
			&report.PurchaseOrdersCount,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	// If specific buyer was requested and not found, return error
	if buyerID != nil && len(reports) == 0 {
		return nil, pkg.ServiceError{
			ResponseCode:  pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode,
			InternalError: pkg.ServiceErrors[pkg.ErrNotFound].InternalError,
		}
	}

	return reports, nil
}
