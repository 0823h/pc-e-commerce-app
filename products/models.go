package products

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"
	"tmdt-backend/common"
	"tmdt-backend/manufacturers"
	"tmdt-backend/users"

	"github.com/elastic/go-elasticsearch/esapi"
)

type Product struct {
	ID             uint64  `gorm:"primaryKey"`
	SKU            string  `gorm:"column:sku"`
	Name           string  `gorm:"column:name"`
	Description    string  `gorm:"column:description"`
	Images         *string `gorm:"column:images"`
	Rating         float32 `gorm:"column:rating"`
	Price          float64 `gorm:"column:price"`
	Quantity       uint    `gorm:"column:quantity"`
	SoldAmount     uint    `gorm:"column:sold_amount"`
	ManufacturerID uint
	Manufacturer   manufacturers.Manufacturer `gorm:"foreignKey:ManufacturerID"`
	IsDeleted      bool                       `gorm:"column:is_deleted;default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func AutoMigrate() {
	db := common.GetDB()

	db.AutoMigrate(&Product{})
	db.AutoMigrate(&Rating{})
}

func SaveOne(data interface{}) error {
	// Save to database
	db := common.GetDB()
	err := db.Save(data).Error

	// Index to ES
	es := common.GetES()

	es_data, es_err := json.Marshal(data)

	if es_err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	req := esapi.IndexRequest{
		Index: "product",
		// DocumentID: "2",
		Body:    bytes.NewReader(es_data),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), "1")
	}
	//Return error
	return err

}

func NewProduct() Product {
	var product Product
	return product
}

// func SaveOneToES() {
// 	es := common.GetES()

// 	data, err := json.Marshal(struct {
// 		Title string `json:"title"`
// 	}{Title: title})
// 	if err != nil {
// 		log.Fatalf("Error marshaling document: %s", err)
// 	}
// }

type Rating struct {
	ID            uint64 `gorm:"primaryKey"`
	UserID        uint64
	User          users.User `gorm:"foreignKey:UserID"`
	ProductID     uint64
	Product       Product `gorm:"foreignKey:ProductID"`
	Rate          uint    `gorm:"column:rate;default:0"`
	NumberOfClick uint    `gorm:"column:number_of_click;default:0"`
	IsDeleted     bool    `gorm:"default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewRating() Rating {
	var rating Rating
	return rating
}
