package manufacturers

import (
	"time"
	"tmdt-backend/common"

	"github.com/google/uuid"
)

type Manufacturer struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string    `gorm:"column:name"`
	Origin    string    `gorm:"column:origin"`
	IsDeleted bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func AutoMigrate() {
	db := common.GetDB()

	db.AutoMigrate(&Manufacturer{})
}

func SaveOne(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}
