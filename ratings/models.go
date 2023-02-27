package ratings

import (
	"tmdt-backend/products"

	"github.com/gin-gonic/gin"
)

func RatingRouter(router *gin.RouterGroup) {
	router.GET("/", products.GetRatings)
}
