package main

import (
	"log"
	"strings"
	"tmdt-backend/common"
	"tmdt-backend/products"
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
	es := common.ESInit()
	index := "products"
	mapping := `{
		"settings": {
			"number_of_shards": 1
		},
		"mappings": {
			"properties": {
				"field1": {
					"type": "text"
				}
			}
		}
	}`

	res, err := es.Indices.Create(
		index,
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	log.Println(res)

	r := gin.Default()
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("GET")

	v1 := r.Group("/api")
	users.UsersRegister(v1.Group("/users"))
	products.ProductRouter(v1.Group("/products"))

	testAuth := r.Group("/api/ping")

	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()

	//fmt.Println(viper.Get("DATABASE_HOST"))
}
