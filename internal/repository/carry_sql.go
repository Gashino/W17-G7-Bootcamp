package repository

import (
	"app/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

	fmt.Println(v)
	_, err = r.db.Exec(
		"INSERT INTO carries (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)",
		v.Cid, v.CompanyName, v.Address, v.Telephone, v.LocalityId,
	)

	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("SQL Error")
		return
	}
	c = v

	return
}

func (r *CarrySql) SearchByLocality(locality_id int) (v map[string]models.CarryByLocality, err error) {
	v = make(map[string]models.CarryByLocality)

	var rows *sql.Rows

	query_str := "SELECT l.id AS locality_id, l.locality_name, COUNT(c.id) AS carries_count FROM localities l LEFT JOIN carries c ON l.id = c.locality_id !WHERE_COND! GROUP BY l.id, l.locality_name"

	if locality_id >= 0 {
		query_str = strings.Replace(query_str, "!WHERE_COND!", "where l.id = ?", -1)
		rows, err = r.db.Query(query_str, locality_id)
	} else {
		query_str = strings.Replace(query_str, "!WHERE_COND!", "", -1)
		rows, err = r.db.Query(query_str)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var carryByLoc models.CarryByLocality
		err := rows.Scan(&carryByLoc.LocalityId, &carryByLoc.LocalityName, &carryByLoc.CarriesCount)
		if err != nil {
			return nil, err
		}
		v[carryByLoc.LocalityId] = carryByLoc
	}

	return
}
