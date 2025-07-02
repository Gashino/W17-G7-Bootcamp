package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type EmployeeHandler struct {
	service service.EmployeeService
}

func NewEmployeeHandler(service service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
	}
}

func (h *EmployeeHandler) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.service.FindAll()
	if err != nil {
		handleServiceError(w, err)
		return
	}

	responses := models.ToEmployeeResponsesFromMap(employees)
	writeResponse(w, http.StatusOK, models.EmployeesResponse{Data: responses}, nil)
}

func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	id, valid := extractIDFromURL(w, r)
	if !valid {
		return
	}

	employee, err := h.service.FindById(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response := employee.ToEmployeeDataResponse()
	writeResponse(w, http.StatusOK, response, nil)
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee models.EmployeeDTO
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceError{
			Code:         400,
			ResponseCode: http.StatusBadRequest,
			Message:      "Invalid request body",
		})
		return
	}

	employeeModel := employee.ToEmployee()

	if err := validateRequest(employeeModel, false); err != nil {
		handleServiceError(w, err)
		return
	}

	newEmployee, err := h.service.Save(employeeModel)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response := newEmployee.ToEmployeeDataResponse()
	writeResponse(w, http.StatusCreated, response, nil)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id, valid := extractIDFromURL(w, r)
	if !valid {
		return
	}

	var employee models.EmployeeDTO
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceError{
			Code:         404,
			ResponseCode: http.StatusBadRequest,
			Message:      "Invalid request body",
		})
		return
	}

	if !validateRequestUpdatePatch(employee) {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceErrors[pkg.ErrBadRequest])
		return
	}

	employeeModel := employee.ToEmployee()
	updatedEmployee, err := h.service.Update(employeeModel, id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response := updatedEmployee.ToEmployeeDataResponse()
	writeResponse(w, http.StatusOK, response, nil)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id, valid := extractIDFromURL(w, r)
	if !valid {
		return
	}

	err := h.service.Delete(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeResponse(w, http.StatusNoContent, nil, nil)
}

// Functions to handle the responses

// writeResponse is a helper function to write JSON responses with status code
func writeResponse(w http.ResponseWriter, status int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err != nil {
		response := struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// handleServiceError is a helper function to handle service errors consistently
func handleServiceError(w http.ResponseWriter, err error) {
	var serviceErr pkg.ServiceError
	if errors.As(err, &serviceErr) {
		writeResponse(w, serviceErr.ResponseCode, nil, err)
	} else {
		writeResponse(w, http.StatusInternalServerError, nil, pkg.ServiceErrors[pkg.ErrInternalServer])
	}
}

// Function to extract the ID from the URL
func extractIDFromURL(w http.ResponseWriter, r *http.Request) (int, bool) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceError{
			Code:         100,
			ResponseCode: http.StatusBadRequest,
			Message:      "ID is required",
		})
		return 0, false
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceError{
			Code:         101,
			ResponseCode: http.StatusBadRequest,
			Message:      "Invalid ID format",
		})
		return 0, false
	}

	return id, true
}

// validateRequest is a helper function to validate the employee request
func validateRequest(employee models.Employee, validateID bool) error {
	if err := models.ValidateEmployee(employee, validateID); err != nil {
		errorR := pkg.ServiceErrors[pkg.ErrBadRequest]
		errorR.Message = err.Error()
		return errorR

	}
	return nil
}

// validateRequestUpdatePatch is a helper function to validate the employee request for update and patch
func validateRequestUpdatePatch(employee models.EmployeeDTO) bool {
	if employee.CardNumberID == "" && employee.FirstName == "" && employee.LastName == "" && employee.WarehouseID == 0 {
		return false
	}
	return true
}
