package products

import "github.com/gin-gonic/gin"

func ProductRouter(router *gin.RouterGroup) {
	router.GET("/", GetAllProducts)
}
