package users

import (
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

func UpdateContextUserModel(c *gin.Context, my_user_id string) {
	var myUserModel User
	if my_user_id != "" {
		db := common.GetDB()
		db.First(&myUserModel, my_user_id)
	}
	c.Set("my_user_id", my_user_id)
	c.Set("my_user_model", myUserModel)
}
