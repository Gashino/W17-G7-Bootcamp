package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"errors"
	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"net/http"
)

func NewProductRecordDefault(sv service.IProductRecordService) *ProductRecordhDefault {
	return &ProductRecordhDefault{sv: sv}
}

type ProductRecordhDefault struct {
	sv service.IProductRecordService
}

func (d ProductRecordhDefault) Create() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		var prodRecord models.ProductRecord

		errParsing := request.JSON(req, &prodRecord)
		if errParsing != nil {
			response.Error(writer, http.StatusInternalServerError, errParsing.Error())
			return
		}

		if isValid := prodRecord.Validate(); !isValid {
			response.Error(writer, http.StatusUnprocessableEntity, "some field is not valid for product_record")
			return
		}

		result, err := d.sv.Create(prodRecord)

		if err != nil {
			svcErr := pkg.ServiceError{}
			if errors.As(err, &svcErr) {
				response.Error(writer, svcErr.ResponseCode, svcErr.Error())
				return
			}
			response.Error(writer, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, pkg.ServiceErrors[pkg.ErrInternalServer].Error())
			return
		}

		response.JSON(writer, http.StatusCreated, map[string]any{
			"data": result,
		})

	}
}
