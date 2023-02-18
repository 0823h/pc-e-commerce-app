package products

import (
	"context"
	"time"
	"tmdt-backend/common"
	"tmdt-backend/manufacturers"

	"github.com/elastic/go-elasticsearch/v8/esapi"
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

	//Return error
	return err

}

func SaveOneToES() {
	es := common.GetES()
	res, err := esapi.IndexRequest("index_name").
		Raw([]byte(`{
	  "id": 1,
	  "name": "Foo",
	  "price": 10
	}`)).Do(context.Background())
}
