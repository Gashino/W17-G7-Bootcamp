package models

import (
	"fmt"
	"regexp"
	"strings"
)

// Helper functions for creating pointers
func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

// Helper functions to check if pointers have valid values
func IsValidStringPtr(s *string) bool {
	return s != nil && strings.TrimSpace(*s) != ""
}

func IsValidIntPtr(i *int) bool {
	return i != nil && *i > 0
}

// Employee represents the employee entity in the database
// This struct will be used as a map in the database
// swagger:model Employee
type Employee struct {
	ID           *int    `json:"id"`
	CardNumberID *string `json:"card_number_id"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	WarehouseID  *int    `json:"warehouse_id"`
}

type EmployeeReport struct {
	ID                 int    `json:"id"`
	CardNumberID       string `json:"card_number_id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	WarehouseID        int    `json:"warehouse_id"`
	InboundOrdersCount int    `json:"inbound_orders_count"`
}

type EmployeeReportsResponse struct {
	Data []EmployeeReport `json:"data"`
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
	// Validate ID if required
	if validateID {
		if employee.ID == nil || *employee.ID <= 0 {
			return fmt.Errorf("id is required and must be greater than 0")
		}
	}

	// Validate CardNumberID
	if employee.CardNumberID == nil {
		return fmt.Errorf("card_number_id is required")
	}
	if len(*employee.CardNumberID) != 8 {
		return fmt.Errorf("card_number_id must be 8 digits")
	}
	cardNumberRegex := regexp.MustCompile(`^[0-9]{8}$`)
	if !cardNumberRegex.MatchString(*employee.CardNumberID) {
		return fmt.Errorf("card_number_id must contain only numbers")
	}

	// Validate FirstName
	if employee.FirstName == nil || len(strings.TrimSpace(*employee.FirstName)) == 0 {
		return fmt.Errorf("first_name is required")
	}

	// Validate LastName
	if employee.LastName == nil || len(strings.TrimSpace(*employee.LastName)) == 0 {
		return fmt.Errorf("last_name is required")
	}

	// Validate WarehouseID
	if employee.WarehouseID == nil || *employee.WarehouseID <= 0 {
		return fmt.Errorf("warehouse_id is required and must be greater than 0")
	}

	return nil
}

// ValidateEmployeeUpdate validates an employee for update operations (allows partial data)
func ValidateEmployeeUpdate(employee Employee) error {
	// At least one field must be provided for update
	if employee.CardNumberID == nil && employee.FirstName == nil &&
		employee.LastName == nil && employee.WarehouseID == nil {
		return fmt.Errorf("at least one field is required for update")
	}

	// Validate CardNumberID if provided
	if employee.CardNumberID != nil {
		if len(*employee.CardNumberID) != 8 {
			return fmt.Errorf("card_number_id must be 8 digits")
		}
		cardNumberRegex := regexp.MustCompile(`^[0-9]{8}$`)
		if !cardNumberRegex.MatchString(*employee.CardNumberID) {
			return fmt.Errorf("card_number_id must contain only numbers")
		}
	}

	// Validate FirstName if provided
	if employee.FirstName != nil && len(strings.TrimSpace(*employee.FirstName)) == 0 {
		return fmt.Errorf("first_name cannot be empty")
	}

	// Validate LastName if provided
	if employee.LastName != nil && len(strings.TrimSpace(*employee.LastName)) == 0 {
		return fmt.Errorf("last_name cannot be empty")
	}

	// Validate WarehouseID if provided
	if employee.WarehouseID != nil && *employee.WarehouseID <= 0 {
		return fmt.Errorf("warehouse_id must be greater than 0")
	}

	return nil
}

// ToEmployeeResponse converts an Employee entity to EmployeeResponse
func (e Employee) ToEmployeeResponse() EmployeeResponse {
	response := EmployeeResponse{}

	if e.ID != nil {
		response.ID = *e.ID
	}
	if e.CardNumberID != nil {
		response.CardNumberID = *e.CardNumberID
	}
	if e.FirstName != nil {
		response.FirstName = *e.FirstName
	}
	if e.LastName != nil {
		response.LastName = *e.LastName
	}
	if e.WarehouseID != nil {
		response.WarehouseID = *e.WarehouseID
	}

	return response
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
		CardNumberID: StringPtr(dto.CardNumberID),
		FirstName:    StringPtr(dto.FirstName),
		LastName:     StringPtr(dto.LastName),
		WarehouseID:  IntPtr(dto.WarehouseID),
	}
}

// ToEmployeeForUpdate converts an EmployeeDTO to Employee entity for update operations
// Only converts non-empty fields to pointers, leaving empty fields as nil
func (dto EmployeeDTO) ToEmployeeForUpdate() Employee {
	employee := Employee{}

	// Only set CardNumberID if it's not empty
	if dto.CardNumberID != "" {
		employee.CardNumberID = StringPtr(dto.CardNumberID)
	}

	// Only set FirstName if it's not empty
	if dto.FirstName != "" {
		employee.FirstName = StringPtr(dto.FirstName)
	}

	// Only set LastName if it's not empty
	if dto.LastName != "" {
		employee.LastName = StringPtr(dto.LastName)
	}

	// Only set WarehouseID if it's not zero
	if dto.WarehouseID != 0 {
		employee.WarehouseID = IntPtr(dto.WarehouseID)
	}

	return employee
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
