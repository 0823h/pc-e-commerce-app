package manufacturers

import "time"

type ManufacturerModel struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"column:name"`
	Origin    string `gorm:"column:origin"`
	IsDeleted bool   `gorm:"column:is_deleted;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
