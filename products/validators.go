package products

import (
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateProductValidator struct {
	SKU            string         `form:"sku" json:"sku" binding:"required"`
	Name           string         `form:"name" json:"name" binding:"required"`
	Description    string         `form:"description" json:"description" binding:"required"`
	Images         pq.StringArray `form:"images" json:"images" binding:"required"`
	ManufacturerId uint           `form:"manufacturer_id" json:"manufacturer_id" binding:"required"`
	productModel   Product        `json:"-"`
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
	// self.productModel.ID = uuid.New()
	self.productModel.SKU = self.SKU
	self.productModel.Name = self.Name
	self.productModel.Description = self.Description
	if (self.Images) != nil {
		self.productModel.Images = self.Images
	}
	self.productModel.ManufacturerID = self.ManufacturerId

	return nil
}

func NewCreateProductValidatorFillWith(productModel Product) CreateProductValidator {
	createProductValidator := NewCreateProductValidator()
	createProductValidator.SKU = productModel.SKU
	createProductValidator.Name = productModel.Name
	createProductValidator.Description = productModel.Description

	if productModel.Images != nil {
		createProductValidator.Images = productModel.Images
	}

	return createProductValidator
}

type RatingValidator struct {
	Rate uint `form:"rate" json:"rate" binding:"required"`
}

func NewRatingValidator() RatingValidator {
	var ratingValidator RatingValidator
	return ratingValidator
}

func (self *RatingValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	return err
}

type UpdateProductValidator struct {
	Product struct {
		Name        string `form:"name" json:"name"`
		Description string `form:"description" json:"description"`
	} `json:"product"`
	productModel Product `json:"-"`
}

func NewUpdateProductValidator() UpdateProductValidator {
	var updateProductValidator UpdateProductValidator
	return updateProductValidator
}

func (self *UpdateProductValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)

	if self.Product.Description != "" {
		self.productModel.Description = self.Product.Description
	}

	if self.Product.Name != "" {
		self.productModel.Name = self.Product.Name
	}

	return err
}

func NewUpdateProductValidatorFillWith(productModel Product) UpdateProductValidator {
	var updateProductValidator UpdateProductValidator
	updateProductValidator.productModel = productModel
	return updateProductValidator
}
