package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

//type Merchant struct {
//	Id           int64          `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // 数据表id
//	MerchantId   string         `json:"merchant_id"`                          // 商户id
//	UserId       int64          `json:"user_id"`                              // 对应用户的id
//	MerchantName string         `json:"merchant_name"`                        // 商户名
//	Description  string         `json:"description"`                          // 商户描述
//	LogoUrl      string         `json:"logo_url"`                             // 商户Logo图片
//	Address      string         `json:"address"`                              // 商户地址
//	CreatedAt    time.Time      `json:"created_at"`
//	UpdatedAt    time.Time      `json:"updated_at"`
//	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
//}

type Merchant struct {
	Id         int64  `json:"id" gorm:"primaryKey;autoIncrement"` // 数据表id
	MerchantId string `json:"merchant_id"`                        // 商户id

	UserId int64 `json:"user_id" gorm:"not null;index"` // 必须有效
	User   User  `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	MerchantName string `json:"merchant_name"` // 商户名
	Description  string `json:"description"`   // 商户描述
	LogoUrl      string `json:"logo_url"`      // 商户Logo图片
	Address      string `json:"address"`       // 商户地址

	AvgStarsFloat float64 `gorm:"column:avg_stars" json:"-"`
	AvgStars      string  `json:"avg_stars"`   // 平均评分
	StarsCount    int64   `json:"stars_count"` // 评论数

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (Merchant) TableName() string {
	return "a_merchant"
}

func (m *Merchant) AfterFind(tx *gorm.DB) error {
	m.AvgStars = fmt.Sprintf("%.2f", m.AvgStarsFloat)
	return nil
}
