package models

import (
	"gorm.io/gorm"
	"time"
)

type Merchant struct {
	Id           int64          `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // 数据表id
	MerchantId   string         `json:"merchant_id"`                          // 商户id
	UserId       int64          `json:"user_id"`                              // 对应用户的id
	MerchantName string         `json:"merchant_name"`                        // 商户名
	Description  string         `json:"description"`                          // 商户描述
	LogoUrl      string         `json:"logo_url"`                             // 商户Logo图片
	Address      string         `json:"address"`                              // 商户地址
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

func (Merchant) TableName() string {
	return "a_merchant"
}
