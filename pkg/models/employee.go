package models

import (
	"fmt"
	"strings"
)

// Employee represents the employee entity in the database
// This struct will be used as a map in the database
// swagger:model Employee
type Employee struct {
	ID           int
	CardNumberID string
	FirstName    string
	LastName     string
	WarehouseID  int
}

type EmployeeDocument struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}

type EmployeeDTO struct {
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}

type EmployeeUpdateDTO struct {
	EmployeeDTO
	ID int `json:"id"`
}

type EmployeeResponse struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}

type EmployeesResponse struct {
	Data []EmployeeResponse `json:"data"`
}

type EmployeeDataResponse struct {
	Data EmployeeResponse `json:"data"`
}

func ValidateEmployee(employee Employee, validateID bool) error {

	if len(employee.CardNumberID) != 8 {
		return fmt.Errorf("card_number_id must be 8 digits")
	}

	if len(strings.TrimSpace(employee.FirstName)) == 0 {
		return fmt.Errorf("first_name is required")
	}

	if len(strings.TrimSpace(employee.LastName)) == 0 {
		return fmt.Errorf("last_name is required")
	}

	if employee.WarehouseID == 0 {
		return fmt.Errorf("warehouse_id must be a valid number")
	}

	if validateID {
		if employee.ID == 0 {
			return fmt.Errorf("id is required")
		}
	}

	return nil
}
