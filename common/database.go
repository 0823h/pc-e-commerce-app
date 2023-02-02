package common

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct
func Init() *gorm.DB {
	// db, err := gorm.Open("postgres")
	fmt.Println(viper.Get("DATABASE_NAME"))
	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		viper.Get("DATABASE_HOST"),
		viper.Get("DATABASE_USER"),
		viper.Get("DATABASE_PASSWORD"),
		viper.Get("DATABASE_NAME"),
		viper.Get("DATABASE_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	DB = db
	return DB
}

// Using this function to get a connection, you can create your connection pool here
func GetDB() *gorm.DB {
	return DB
}
