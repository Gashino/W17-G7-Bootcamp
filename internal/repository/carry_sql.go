package repository

import (
	"app/pkg/models"
	"database/sql"
	"errors"
	"fmt"
)

// CarrySql is a function that returns a new instance of CarrySql
func NewCarrySql(db *sql.DB) *CarrySql {
	return &CarrySql{db: db}
}

// Struct for CarrySql Repository
type CarrySql struct {
	// db is a map of warehouse
	db *sql.DB
}

func (r *CarrySql) Create(v models.Carry) (c models.Carry, err error) {
	_, err = r.db.Exec(
		"INSERT INTO carries (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)",
		c.Cid, c.CompanyName, c.Address, c.Telephone, c.LocalityId,
	)

	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("SQL Error")
		return
	}
	c = v

	return
}

func (r *CarrySql) SearchByLocality(locality_id int) (v map[int]models.Carry, err error) {
	v = make(map[int]models.Carry)
	query_str := "SELECT id, cid, company_name, address, telephone, locality_id FROM carries"
	var rows *sql.Rows

	if locality_id >= 0 {
		query_str += " WHERE locality_id = ?"
		rows, err = r.db.Query(query_str, locality_id)
	} else {
		rows, err = r.db.Query(query_str)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var carry models.Carry
		err := rows.Scan(&carry.ID, &carry.Cid, &carry.CompanyName, &carry.Address, &carry.Telephone, &carry.LocalityId)
		if err != nil {
			return nil, err
		}
		v[carry.ID] = carry
	}

	return
}
