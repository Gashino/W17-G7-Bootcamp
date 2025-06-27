package models

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
