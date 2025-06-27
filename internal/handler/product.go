package handler

import (
	"app/internal/service"
	"app/pkg/models"
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

		dataResponse := make(map[int]models.ProductDoc)
		for key, value := range result {
			dataResponse[key] = value.ToJSON()
		}

		response.JSON(writer, http.StatusOK, dataResponse)
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

		response.JSON(writer, http.StatusOK, result.ToJSON())
	}
}

func (d ProductDefault) Create() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (d ProductDefault) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (d ProductDefault) Patch() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
