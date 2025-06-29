package models

import (
	"gorm.io/gorm"
	"time"
)

type Goods struct {
	Id          int64          `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // 数据表id
	MerchantId  int64          `json:"merchant_id"`                          // 对应商户id
	Description string         `json:"description"`                          // 商品描述
	LogoUrl     string         `json:"logo_url"`                             // 商品Logo图片
	Price       float32        `json:"price"`                                // 价格
	Residue     int64          `json:"residue"`                              // 剩余库存
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"Updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (Goods) TableName() string {
	return "a_goods"
}
