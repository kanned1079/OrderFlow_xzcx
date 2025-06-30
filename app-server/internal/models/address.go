package models

import (
	"gorm.io/gorm"
	"time"
)

type Address struct {
	Id          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      int64          `gorm:"index" json:"user_id"`
	FullName    string         `json:"full_name"`
	PhoneNumber string         `json:"phone_number"`
	FullAddress string         `json:"full_address"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (Address) TableName() string {
	return "a_address"
}
