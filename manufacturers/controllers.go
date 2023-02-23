package manufacturers

import (
	"net/http"
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

func CreateManufacturer(c *gin.Context) {
	validator := NewCreateManufacturerValidator()

	if err := validator.Bind(c); err != nil {
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	if err := SaveOne(&validator.manufacturerModel); err != nil {
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	serializer := ManufacturerSerializer{c, validator.manufacturerModel}
	common.SendResponse(c, http.StatusCreated, "Success", serializer.Response())
	return
}
