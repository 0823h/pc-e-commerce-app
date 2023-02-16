package products

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"reflect"
	"time"
	"tmdt-backend/common"
	"tmdt-backend/manufacturers"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/google/uuid"
)

type Product struct {
	ID             uuid.UUID `gorm:"primaryKey"`
	SKU            string    `gorm:"column:sku"`
	Name           string    `gorm:"column:name"`
	Description    string    `gorm:"column:description"`
	Images         *string   `gorm:"column:images"`
	Rating         float32   `gorm:"column:rating"`
	Price          float64   `gorm:"column:price"`
	Quantity       uint      `gorm:"column:quantity"`
	SoldAmount     uint      `gorm:"column:sold_amount"`
	ManufacturerID uint
	Manufacturer   manufacturers.Manufacturer `gorm:"foreignKey:ManufacturerID"`
	IsDeleted      bool                       `gorm:"column:is_deleted;default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func AutoMigrate() {
	db := common.GetDB()

	db.AutoMigrate(&Product{})
}

func SaveOne(data interface{}) error {
	// Save to database
	db := common.GetDB()
	err := db.Save(data).Error

	// Save to Elasticsearch
	// Build the request body.
	es := common.GetES()
	es_data, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	es_data_id := reflect.ValueOf(data).Elem().FieldByName("ID").String()

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "product",
		DocumentID: es_data_id,
		Body:       bytes.NewReader(es_data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), es_data_id)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	defer res.Body.Close()

	//Return error
	return err

}

// func SaveOne(data interface{}) error {
// 	db := common.GetDB()
// 	err := db.Save(data).Error
// 	return err
// }
