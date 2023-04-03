package categories

import (
	"net/http"
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	validator := NewCreateCategoryValidator()
	if err := validator.Bind(c); err != nil {
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	db := common.GetDB()
	if err := db.Create(&validator.categoryModel).Error; err != nil {
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	common.SendResponse(c, http.StatusCreated, "Success", validator.categoryModel)
	return
}

func GetCategories(c *gin.Context) {
	db := common.GetDB()
	var categories []Category

	db.Find(&categories)
	common.SendResponse(c, http.StatusUnprocessableEntity, "Success", categories)
	return
}
