package category_product_relation

import (
	"time"
	"tmdt-backend/categories"
	"tmdt-backend/products"
)

type CategoryProductRelation struct {
	ID         uint64 `gorm:"primaryKey"`
	ProductID  uint
	Product    products.Product `gorm:"foreignKey:ProductID"`
	CategoryID uint
	Category   categories.Category `gorm:"foreignKey:CategoryID"`
	IsDeleted  bool                `gorm:"column:is_deleted;default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
