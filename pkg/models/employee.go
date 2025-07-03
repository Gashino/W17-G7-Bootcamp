package models

import (
	"fmt"
	"regexp"
	"strings"
)

// Employee represents the employee entity in the database
// This struct will be used as a map in the database
// swagger:model Employee
type Employee struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
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

	// Validate that CardNumberID contains only digits
	cardNumberRegex := regexp.MustCompile(`^[0-9]{8}$`)
	if !cardNumberRegex.MatchString(employee.CardNumberID) {
		return fmt.Errorf("card_number_id must contain only numbers")
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

// ToEmployeeResponse converts an Employee entity to EmployeeResponse
func (e Employee) ToEmployeeResponse() EmployeeResponse {
	return EmployeeResponse{
		ID:           e.ID,
		CardNumberID: e.CardNumberID,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		WarehouseID:  e.WarehouseID,
	}
}

// ToEmployeeDataResponse converts an Employee entity to EmployeeDataResponse
func (e Employee) ToEmployeeDataResponse() EmployeeDataResponse {
	return EmployeeDataResponse{
		Data: e.ToEmployeeResponse(),
	}
}

// ToEmployee converts an EmployeeDTO to Employee entity
func (dto EmployeeDTO) ToEmployee() Employee {
	return Employee{
		CardNumberID: dto.CardNumberID,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		WarehouseID:  dto.WarehouseID,
	}
}

// ToEmployeeResponses converts a slice of Employee entities to a slice of EmployeeResponse
func ToEmployeeResponses(employees []Employee) []EmployeeResponse {
	responses := make([]EmployeeResponse, len(employees))
	for i, e := range employees {
		responses[i] = e.ToEmployeeResponse()
	}
	return responses
}

// ToEmployeeResponsesFromMap converts a map of Employee entities to a slice of EmployeeResponse
func ToEmployeeResponsesFromMap(employees map[int]Employee) []EmployeeResponse {
	responses := make([]EmployeeResponse, 0, len(employees))
	for _, e := range employees {
		responses = append(responses, e.ToEmployeeResponse())
	}
	return responses
}
