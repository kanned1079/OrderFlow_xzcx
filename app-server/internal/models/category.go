package models

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	ID         int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantId int64  `gorm:"index" json:"merchant_id"`
	Title      string `json:"title"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (Category) TableName() string {
	return "a_category"
}
