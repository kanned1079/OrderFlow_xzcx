package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id          int64          `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // 数据表id
	PlatformId  string         `json:"platform_id"`                          // 美团id
	AppOpenId   string         `json:"app_open_id" gorm:"unique"`            // 小程序appId
	SessionKey  string         `json:"session_key"`
	PhoneNumber string         `json:"phone_number"`
	Name        string         `json:"name"`
	Role        string         `json:"role"` // 用户v0 站长v1 管理v10
	LastLoginAt time.Time      `json:"last_login_at" gorm:"type:timestamp"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt   time.Time      `json:"Updated_at" gorm:"type:timestamp"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
}

func (User) TableName() string {
	return "s_user"
}
