package products

import (
	"net/http"
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	var products []Product
	pagination := common.NewPagination()

	db := common.GetDB()
	db.Scopes(common.Paginate(products, &pagination, db)).Joins("Manufacturer").Find(&products)
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func CreateProduct(c *gin.Context) {
	validator := NewCreateProductValidator()

	if err := validator.Bind(c); err != nil {
		// c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := SaveOne(&validator.productModel); err != nil {
		// c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	serializer := ProductSerializer{c, validator.productModel}
	c.JSON(http.StatusCreated, gin.H{"Product": serializer.Response()})
}
