package recommendation

import "github.com/gin-gonic/gin"

func RecommendationRouter(router *gin.RouterGroup) {
	router.GET("/recommendation/:user_id")
}
