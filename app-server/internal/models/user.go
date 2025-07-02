package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id          int64      `json:"id" gorm:"primary_key;AUTO_INCREMENT"` // 数据表id
	PhoneNumber string     `json:"phone_number"`                         // 手机号
	Username    string     `json:"username"`                             // 用户名
	Role        string     `json:"role"`                                 // 用户user 商家trader 管理admin
	Password    string     `json:"password"`                             // 用户密码
	AvatarUrl   string     `json:"avatar_url"`                           // 用户头像URL
	Status      bool       `json:"status"`                               // 账户状态
	LastLoginAt *time.Time `json:"last_login_at"`                        // 上一次登录的时间

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (User) TableName() string {
	return "a_user"
}
