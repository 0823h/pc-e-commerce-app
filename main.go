package main

import (
	"tmdt-backend/common"
	"tmdt-backend/products"
	"tmdt-backend/ratings"
	"tmdt-backend/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	users.AutoMigrate()
	products.AutoMigrate()
}

func main() {
	// Load configuration
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	//Database
	db := common.Init()
	Migrate(db)

	//Elasticsearch
	common.ESInit()

	r := gin.Default()
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("GET")

	// Gorse
	// gorse := client.NewGorseClient("http://127.0.0.1:8087", "api_key")
	// common.GorseInit()
	// gorse := common.GetGorse()
	// fmt.Println(gorse)

	v1 := r.Group("/api")
	users.UsersRegister(v1.Group("/users"))
	products.ProductRouter(v1.Group("/products"))
	ratings.RatingRouter(v1.Group("/ratings"))

	testAuth := r.Group("/api/ping")

	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()

	//fmt.Println(viper.Get("DATABASE_HOST"))
}
