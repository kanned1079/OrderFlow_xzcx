package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	Id          int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderId     string          `gorm:"index" json:"order_id"`    // 所属订单
	UserId      int64           `gorm:"index" json:"user_id"`     // 冗余，方便查询
	MerchantId  int64           `gorm:"index" json:"merchant_id"` // 冗余，方便统计/过滤
	CommentText string          `gorm:"type:text" json:"comment_text"`
	ImagesUrls  json.RawMessage `gorm:"type:json" json:"images_urls"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (Comment) TableName() string {
	return "a_comment"
}
