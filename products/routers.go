package products

import (
	"tmdt-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ProductRouter(router *gin.RouterGroup) {
	router.GET("/", GetAllProducts)
	router.POST("/", CreateProduct)
	//ElasticSearch
	router.GET("/es", GetAllProductsES)
	router.POST("/:id/rate", middlewares.Authorization, RateProduct)
}
