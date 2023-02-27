package middlewares

import (
	"net/http"
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

func Authorization(c *gin.Context) {
	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		common.SendResponse(c, http.StatusUnauthorized, "token not found!", "")
		c.Abort()
		return
	}
	claims, err := common.ValidateToken(clientToken)
	if err != "" {
		// common.SendResponse(c, http.StatusUnauthorized, "token validate error", "")
		common.SendResponse(c, http.StatusUnauthorized, err, "")
		c.Abort()
		return
	}
	c.Set("id", claims.ID)
	c.Set("email", claims.Email)
	c.Next()
}
