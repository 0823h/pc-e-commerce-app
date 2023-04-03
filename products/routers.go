package products

import (
	"tmdt-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ProductRouter(router *gin.RouterGroup) {
	router.GET("/", GetAllProducts)
	router.POST("/", CreateProduct)
	//ElasticSearch
	// router.GET("/es", GetAllProductsES)
	router.POST("/:id/rate", middlewares.Authorization, RateProduct)
	router.PUT("/:id", UpdateProduct)
	router.POST("/image-test", TestImageUpload)
	router.GET("/category/:id", GetCategoryProduct)
}
