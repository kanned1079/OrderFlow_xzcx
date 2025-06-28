package dto

import "stay-server/internal/models"

type CheckUserIsValidPreRegisterRequestDto struct {
	PlatformId string `json:"platform_id"`
	Name       string `json:"name"`
}

type UserLoginRequestDto struct {
	Code          string `json:"code"`           //wx.login 获取的临时登录凭证
	EncryptedData string `json:"encrypted_data"` // 加密的手机号数据
	Iv            string `json:"iv"`             // 加密向量
}

type UserLoginResponseDto struct {
	User         models.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}
