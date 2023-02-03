package products

import (
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

type CreateProductValidator struct {
	Product struct {
		SKU         string `form:"sku" json:"sku" binding:"requireed"`
		Name        string `form:"name" json:"name" binding:"required"`
		Description string `form:"description" json:"description" binding:"required"`
		Images      string `form:"image" json:"image" binding:"required"`
	} `json:"product"`
	productModel ProductModel `json:"-"`
}

func NewCreateProductValidator() CreateProductValidator {
	createProductValidator := CreateProductValidator{}
	return createProductValidator
}

func (self *CreateProductValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.productModel.SKU = self.Product.SKU
	self.productModel.Name = self.Product.Name
	self.productModel.Description = self.Product.Description

	if self.Product.Images != "" {
		self.productModel.Images = &self.Product.Images
	}
	return nil
}

func NewCreateProductValidatorFillWith(productModel ProductModel) CreateProductValidator {
	createProductValidator := NewCreateProductValidator()
	createProductValidator.Product.SKU = productModel.SKU
	createProductValidator.Product.Name = productModel.Name
	createProductValidator.Product.Description = productModel.Description

	if productModel.Images != nil {
		createProductValidator.Product.Images = *productModel.Images
	}

	return createProductValidator
}
