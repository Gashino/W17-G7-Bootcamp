package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

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

func validateRequest(employee models.Employee, validateID bool) error {
	if err := models.ValidateEmployee(employee, validateID); err != nil {
		errorR := pkg.ServiceErrors[pkg.ErrBadRequest]
		errorR.Message = err.Error()
		return errorR

	}
	return nil
}

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
		writeResponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	var responses []models.EmployeeResponse
	for _, e := range employees {
		responses = append(responses, models.EmployeeResponse{
			ID:           e.ID,
			CardNumberID: e.CardNumberID,
			FirstName:    e.FirstName,
			LastName:     e.LastName,
			WarehouseID:  e.WarehouseID,
		})
	}

	writeResponse(w, http.StatusOK, models.EmployeesResponse{Data: responses}, nil)
}

func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceError{
			Code:         100,
			ResponseCode: http.StatusBadRequest,
			Message:      "ID is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceError{
			Code:         101,
			ResponseCode: http.StatusBadRequest,
			Message:      "Invalid ID format",
		})
		return
	}

	employee, err := h.service.FindById(id)
	if err != nil {
		if err == pkg.ServiceErrors[pkg.ErrNotFound] {
			writeResponse(w, http.StatusNotFound, nil, err)
		} else {
			writeResponse(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response := models.EmployeeDataResponse{
		Data: models.EmployeeResponse{
			ID:           employee.ID,
			CardNumberID: employee.CardNumberID,
			FirstName:    employee.FirstName,
			LastName:     employee.LastName,
			WarehouseID:  employee.WarehouseID,
		},
	}

	writeResponse(w, http.StatusOK, response, nil)
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee models.EmployeeDTO
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, nil, pkg.ServiceError{
			Code:         102,
			ResponseCode: http.StatusUnprocessableEntity,
			Message:      "Invalid request body",
		})
		return
	}

	employeeModel := models.Employee{
		CardNumberID: employee.CardNumberID,
		FirstName:    employee.FirstName,
		LastName:     employee.LastName,
		WarehouseID:  employee.WarehouseID,
	}

	if err := validateRequest(employeeModel, false); err != nil {
		var serviceErr pkg.ServiceError
		if errors.As(err, &serviceErr) {
			writeResponse(w, serviceErr.ResponseCode, nil, err)
		} else {
			writeResponse(w, http.StatusInternalServerError, nil, pkg.ServiceErrors[pkg.ErrInternalServer])
		}
		return
	}

	newEmployee, err := h.service.Save(employeeModel)
	if err != nil {
		if err.Error() == "Card ID already exists" {
			writeResponse(w, http.StatusConflict, nil, err)
		} else {
			writeResponse(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response := models.EmployeeDataResponse{
		Data: models.EmployeeResponse{
			ID:           newEmployee.ID,
			CardNumberID: newEmployee.CardNumberID,
			FirstName:    newEmployee.FirstName,
			LastName:     newEmployee.LastName,
			WarehouseID:  newEmployee.WarehouseID,
		},
	}

	writeResponse(w, http.StatusCreated, response, nil)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceError{
			Code:         404,
			ResponseCode: http.StatusBadRequest,
			Message:      "ID is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceErrors[pkg.ErrBadRequest])
		return
	}

	var employee models.EmployeeUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceError{
			Code:         404,
			ResponseCode: http.StatusBadRequest,
			Message:      "Invalid request body",
		})
		return
	}
	employeeModel := models.Employee{
		CardNumberID: employee.CardNumberID,
		FirstName:    employee.FirstName,
		LastName:     employee.LastName,
		WarehouseID:  employee.WarehouseID,
		ID:           employee.ID,
	}
	if err := validateRequest(employeeModel, true); err != nil {
		var serviceErr pkg.ServiceError
		if errors.As(err, &serviceErr) {
			writeResponse(w, serviceErr.ResponseCode, nil, err)
		} else {
			writeResponse(w, http.StatusInternalServerError, nil, pkg.ServiceErrors[pkg.ErrInternalServer])
		}
		return
	}
	updatedEmployee, err := h.service.Update(employeeModel, id)
	if err != nil {
		if err == pkg.ServiceErrors[pkg.ErrNotFound] {
			writeResponse(w, http.StatusNotFound, nil, err)
		} else {
			writeResponse(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response := models.EmployeeDataResponse{
		Data: models.EmployeeResponse{
			ID:           updatedEmployee.ID,
			CardNumberID: updatedEmployee.CardNumberID,
			FirstName:    updatedEmployee.FirstName,
			LastName:     updatedEmployee.LastName,
			WarehouseID:  updatedEmployee.WarehouseID,
		},
	}

	writeResponse(w, http.StatusOK, response, nil)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceError{
			Code:         106,
			ResponseCode: http.StatusBadRequest,
			Message:      "ID is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil, pkg.ServiceError{
			Code:         107,
			ResponseCode: http.StatusBadRequest,
			Message:      "Invalid ID format",
		})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		if err == pkg.ServiceErrors[pkg.ErrNotFound] {
			writeResponse(w, http.StatusNotFound, nil, err)
		} else {
			writeResponse(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	writeResponse(w, http.StatusNoContent, nil, nil)
}
