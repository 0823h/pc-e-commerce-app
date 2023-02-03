package products

import (
	"time"
	"tmdt-backend/manufacturers"
)

type ProductModel struct {
	ID             uint    `gorm:"primaryKey"`
	SKU            string  `gorm:"column:sku"`
	Name           string  `gorm:"column:name"`
	Description    string  `gorm:"column:description"`
	Images         *string `gorm:"column:images"`
	Rating         float32 `gorm:"column:rating"`
	Price          float64 `gorm:"column:price"`
	Quantity       uint    `gorm:"column:quantity"`
	SoldAmount     uint    `gorm:"column:sold_amount"`
	ManufacturerId uint
	Manufacturer   manufacturers.ManufacturerModel `gorm:"foreignKey:manufacturer_id"`
	IsDeleted      bool                            `gorm:"column:is_deleted;default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
