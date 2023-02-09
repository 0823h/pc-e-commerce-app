package products

import (
	"net/http"
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	var products []Product
	pagination := common.NewPagination()
	common.GetPaginationParameter(c, &pagination)

	db := common.GetDB()
	db.Scopes(common.Paginate(products, &pagination, db)).Joins("Manufacturer").Find(&products)
	serializer := ProductsSerializer{c, products}
	pagination.Data = serializer.Response()
	// c.JSON(http.StatusOK, gin.H{"products": products})
	common.SendResponse(c, http.StatusOK, "Success", pagination)
	return
}

func CreateProduct(c *gin.Context) {
	validator := NewCreateProductValidator()

	if err := validator.Bind(c); err != nil {
		// c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		// c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	if err := SaveOne(&validator.productModel); err != nil {
		// c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		// c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	serializer := ProductSerializer{c, validator.productModel}
	// c.JSON(http.StatusCreated, gin.H{"Product": serializer.Response()})
	common.SendResponse(c, http.StatusCreated, "Success", serializer.Response())
	return
}

func UpdateProduct(c *gin.Context) {

}
