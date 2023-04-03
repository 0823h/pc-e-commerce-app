package categories

import (
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

type CreateCategoryValidator struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	categoryModel Category `json:"-"`
}

func NewCreateCategoryValidator() CreateCategoryValidator {
	var validator CreateCategoryValidator
	return validator
}

func (self *CreateCategoryValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.categoryModel.Name = self.Name
	self.categoryModel.Description = self.Description
	return nil
}
