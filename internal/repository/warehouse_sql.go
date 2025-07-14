package repository

import (
	"app/pkg/models"
	"database/sql"
	"errors"
	"fmt"
)

// NewVehicleMap is a function that returns a new instance of VehicleMap
func NewWarehouseSql(db *sql.DB) *WarehouseSql {
	return &WarehouseSql{db: db}
}

// Struct for Warehouse Repository
type WarehouseSql struct {
	// db is a map of warehouse
	db *sql.DB
}

func (r *WarehouseSql) FindAll() (v map[int]models.Warehouse, err error) {
	v = make(map[int]models.Warehouse)

	rows, err := r.db.Query("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouse models.Warehouse
		err := rows.Scan(&warehouse.ID, &warehouse.WarehouseCode, &warehouse.Address, &warehouse.Telephone, &warehouse.MinCapacity, &warehouse.MinTemperature)
		if err != nil {
			return nil, err
		}
		v[warehouse.ID] = warehouse
	}

	return
}

func (r *WarehouseSql) FindByID(id int) (v models.Warehouse, err error) {

	rows, err := r.db.Query("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses WHERE id = ?", id)
	if err != nil {
		err = errors.New("warehouse not found")
		return
	}
	defer rows.Close()

	if !rows.Next() {
		err = errors.New("warehouse not found")
		return
	}

	var warehouse models.Warehouse
	err = rows.Scan(&warehouse.ID, &warehouse.WarehouseCode, &warehouse.Address, &warehouse.Telephone, &warehouse.MinCapacity, &warehouse.MinTemperature)
	if err != nil {
		return models.Warehouse{}, err
	}
	v = warehouse
	return
}

func (r *WarehouseSql) Add(v models.Warehouse) (w models.Warehouse, err error) {

	_, err = r.db.Exec(
		"INSERT INTO warehouses (warehouse_code, address, telephone, minimun_capacity, minimun_temperature) VALUES (?, ?, ?, ?, ?)",
		v.WarehouseCode, v.Address, v.Telephone, v.MinCapacity, v.MinTemperature,
	)

	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("SQL Error")
		return
	}
	w = v

	return
}

func (r *WarehouseSql) FindWarehouseByCode(code string) (v models.Warehouse, err error) {
	rows, err := r.db.Query("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses WHERE warehouse_code = ?", code)
	if err != nil {
		err = errors.New("warehouse not found")
		return
	}
	defer rows.Close()

	if !rows.Next() {
		err = errors.New("warehouse not found")
		return
	}

	var warehouse models.Warehouse
	err = rows.Scan(&warehouse.ID, &warehouse.WarehouseCode, &warehouse.Address, &warehouse.Telephone, &warehouse.MinCapacity, &warehouse.MinTemperature)
	if err != nil {
		return models.Warehouse{}, err
	}

	return
}

func (r *WarehouseSql) FindAvailableID() (id int, err error) {
	return
}

func (r *WarehouseSql) Update(id int, v models.Warehouse) (w models.Warehouse, err error) {
	_, err = r.db.Exec("UPDATE warehouses SET warehouse_code = ?, address = ?, telephone = ?, minimun_capacity= ?, minimun_temperature = ? WHERE id = ?",
		v.WarehouseCode, v.Address, v.Telephone, v.MinCapacity, v.MinTemperature, id)
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("SQL Error")
		return
	}
	w = v
	return
}

func (r *WarehouseSql) Delete(id int) (err error) {

	_, err = r.db.Exec("DELETE FROM warehouses WHERE id = ?", id)
	if err != nil {
		err = errors.New("SQL Error")
	}

	return nil
}
