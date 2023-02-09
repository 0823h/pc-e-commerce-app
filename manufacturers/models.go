package manufacturers

import (
	"time"
	"tmdt-backend/common"
)

type ManufacturerModel struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"column:name"`
	Origin    string `gorm:"column:origin"`
	IsDeleted bool   `gorm:"column:is_deleted;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func AutoMigrate() {
	db := common.GetDB()

	db.AutoMigrate(&ManufacturerModel{})
}

func SaveOne(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}
