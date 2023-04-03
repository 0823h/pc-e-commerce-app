package categories

import (
	"github.com/gin-gonic/gin"
)

func CategoryRouter(router *gin.RouterGroup) {
	router.GET("/", GetCategories)
	router.POST("/", CreateCategory)
	// router.GET("/:id/products", products.GetCategoryProduct)
}
