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
	r.Run()

	fmt.Println(viper.Get("DATABASE_HOST"))
}
