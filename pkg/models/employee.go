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

// EmployeeDTO is the DTO for creating a new employee
// swagger:model EmployeeDTO
type EmployeeDTO struct {
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  string `json:"warehouse_id"`
}

// EmployeeResponse is the DTO for responding with employee data
// swagger:model EmployeeResponse
type EmployeeResponse struct {
	ID           string `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  string `json:"warehouse_id"`
}

// EmployeesResponse is the DTO for responding with multiple employees
type EmployeesResponse struct {
	Data []EmployeeResponse `json:"data"`
}

// EmployeeDataResponse is the DTO for responding with a single employee
type EmployeeDataResponse struct {
	Data EmployeeResponse `json:"data"`
}
