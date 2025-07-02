package models

import (
	"gorm.io/gorm"
	"time"
)

//type Order struct {
//	OrderId       string         `json:"order_id" gorm:"primary_key;"` // 后端服务器生成
//	UserId        int64          `json:"user_id"`                      // 对应用户的id
//	MerchantId    int64          `json:"merchant_id"`                  // 商户id
//	AddressId     int64          `json:"address_id"`                   // 用户地址信息的id
//	Status        string         `json:"status"`                       // pending_accept processing completed_unreviewed completed_reviewed cancelled
//	FailureReason string         `json:"failure_reason"`               // 订单失败原因
//	CreatedAt     time.Time      `json:"created_at"`
//	UpdatedAt     time.Time      `json:"Updated_at"`
//	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
//}

type Order struct {
	OrderId string `json:"order_id" gorm:"primaryKey"`
	UserId  int64  `json:"user_id" gorm:"not null;index"`
	User    User   `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`

	MerchantId int64 `json:"merchant_id" gorm:"not null;index"`
	Merchant   User  `gorm:"foreignKey:MerchantId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`

	AddressId int64   `json:"address_id" gorm:"not null"`
	Address   Address `gorm:"foreignKey:AddressId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`

	Status        string `json:"status"`
	FailureReason string `json:"failure_reason"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type OrderItem struct {
	Id         int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderId    string  `gorm:"index" json:"order_id"`    // 所属订单
	UserId     int64   `gorm:"index" json:"user_id"`     // 冗余，方便查询
	MerchantId int64   `gorm:"index" json:"merchant_id"` // 冗余，方便统计/过滤
	GoodsId    int64   `gorm:"index" json:"goods_id"`    // 商品ID
	GoodsName  string  `json:"goods_name"`               // 商品快照（下单时保存）
	LogoUrl    string  `json:"logo_url"`
	Price      float32 `json:"price"`
	Quantity   int64   `json:"quantity"`
	Total      float32 `json:"total"` // Price * Quantity

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (Order) TableName() string {
	return "a_order"
}

func (OrderItem) TableName() string {
	return "a_order_item"
}
