package manufacturers

import (
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

type CreateManufacturerValidator struct {
	Manufacturer struct {
		Name   string `form:"name" json:"name" binding:"required"`
		Origin string `form:"origin" json:"origin" binding:"required"`
	} `json:"manufacturer"`
	manufacturerModel Manufacturer `json:"-"`
}

func NewCreateManufacturerValidator() CreateManufacturerValidator {
	var createManufacturerValidator CreateManufacturerValidator
	return createManufacturerValidator
}

func (self *CreateManufacturerValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	// self.manufacturerModel.ID = uuid.New()
	self.manufacturerModel.Name = self.Manufacturer.Name
	self.manufacturerModel.Origin = self.Manufacturer.Origin

	return err
}
