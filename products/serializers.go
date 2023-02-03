package products

import "github.com/gin-gonic/gin"

type ProductSerializer struct {
	C *gin.Context
	ProductModel
}

type ProductRe