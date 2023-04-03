package categories

import (
	"time"
	"tmdt-backend/common"
)

type Category struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	IsDeleted   bool   `gorm:"column:is_deleted;default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func AutoMigrate() {
	db := common.GetDB()
	db.AutoMigrate(&Category{})
}

func NewCategory() Category {
	var new_category Category
	return new_category
}
