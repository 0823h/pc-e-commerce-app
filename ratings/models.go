package ratings

import (
	"time"
	"tmdt-backend/products"
	"tmdt-backend/users"
)

type Rating struct {
	ID            string `gorm:"primaryKey"`
	UserID        string
	User          users.User `gorm:"foreignKey:UserID"`
	ProductID     string
	Product       products.Product `gorm:"foreignKey:ProductID"`
	Rating        uint
	NumberOfClick uint
	IsDeleted     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
