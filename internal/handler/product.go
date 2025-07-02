package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"errors"
	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func NewProductDefault(sv service.IProductService) *ProductDefault {
	return &ProductDefault{sv: sv}
}

type ProductDefault struct {
	sv service.IProductService
}

func (d ProductDefault) GetAll() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		result, err := d.sv.GetAll()

		if err != nil {
			response.Error(writer, http.StatusInternalServerError, err.Error())
			return
		}

		response.JSON(writer, http.StatusOK, map[string]any{
			"message": "success",
			"data":    result,
		})
	}
}

func (d ProductDefault) GetById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(request, "id"))

		result, err := d.sv.GetById(id)

		if err != nil {
			response.Error(writer, http.StatusNotFound, err.Error())
			return
		}

		response.JSON(writer, http.StatusOK, map[string]any{
			"message": "success",
			"data":    result,
		})
	}
}

func (d ProductDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var productDoc models.Product

		errParsing := request.JSON(r, &productDoc)
		if errParsing != nil {
			response.Error(w, http.StatusInternalServerError, errParsing.Error())
			return
		}

		if isValid := productDoc.Validate(); !isValid {
			response.Error(w, http.StatusUnprocessableEntity, "some field is not valid for product")
			return
		}

		result := d.sv.Create(productDoc)

		if result != nil {
			response.Error(w, http.StatusConflict, result.Error())
			return
		}

		response.JSON(w, http.StatusCreated, "product successfully created")
	}
}

func (d ProductDefault) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(request, "id"))

		result := d.sv.Delete(id)

		if result != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(result, &svcErr) {
				response.Error(writer, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(writer, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		response.JSON(writer, http.StatusNoContent, nil)
	}
}

func (d ProductDefault) Patch() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		var productJson models.Product
		id, _ := strconv.Atoi(chi.URLParam(req, "id"))

		err := request.JSON(req, &productJson)
		if err != nil {
			response.Error(writer, http.StatusInternalServerError, err.Error())
			return
		}

		result, errServ := d.sv.Update(id, productJson)

		if errServ != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(errServ, &svcErr) {
				response.Error(writer, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(writer, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return

		}

		response.JSON(writer, http.StatusOK, map[string]any{
			"message": "success",
			"data":    result,
		})
	}
}
