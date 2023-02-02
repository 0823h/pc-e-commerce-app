package main

import (
	"fmt"
	"tmdt-backend/common"
	"tmdt-backend/users"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	users.AutoMigrate()
}

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	db := common.Init()
	Migrate(db)

	r := gin.Default()

	v1 := r.Group("/api")
	users.UsersRegister(v1.Group("/users"))

	testAuth := r.Group("/api/ping")

	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run()

	fmt.Println(viper.Get("DATABASE_HOST"))
}
