package models

import (
	"fmt"
	"strconv"
	"strings"
)

// Employee represents the employee entity in the database
// This struct will be used as a map in the database
// swagger:model Employee
type Employee struct {
	ID           string
	CardNumberID string
	FirstName    string
	LastName     string
	WarehouseID  string
}

type EmployeeDTO struct {
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  string `json:"warehouse_id"`
}

type EmployeeUpdateDTO struct {
	EmployeeDTO
	ID string `json:"id"`
}

type EmployeeResponse struct {
	ID           string `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  string `json:"warehouse_id"`
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

	if _, err := strconv.Atoi(employee.WarehouseID); err != nil {
		return fmt.Errorf("warehouse_id must be a valid number")
	}

	if validateID {
		if len(strings.TrimSpace(employee.ID)) == 0 {
			return fmt.Errorf("id is required")
		}
	}

	return nil
}
