package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"stay-server/internal/config"
	"stay-server/internal/models"
	"time"
)

func (this *Utils) GenerateAccessToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Name,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * time.Duration(config.AppCfg.Runtime.AccessTokenExpiredIn)).Unix(), // 6 小时有效期
	})
	return token.SignedString(this.jwtSecret)
}
