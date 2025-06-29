package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"stay-server/internal/config"
	"stay-server/internal/models"
	"strings"
	"time"
)

func (this *Utils) GenerateAccessToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":           user.Id,
		"phone_number": user.PhoneNumber,
		"role":         user.Role,
		"exp":          time.Now().Add(time.Hour * time.Duration(config.AppCfg.Runtime.AccessTokenExpiredIn)).Unix(), // 6 小时有效期
	})
	//this.Logger.PrintInfo("jwt secret: ", config.AppCfg.Runtime.JwtSecret)

	return token.SignedString([]byte(config.AppCfg.Runtime.JwtSecret))
}

func (this *Utils) ExtractTokenClaims(ctx *gin.Context) (jwt.MapClaims, error) {
	authHeader := ctx.GetHeader("Authorization")
	//this.Logger.PrintInfo(authHeader)
	if authHeader == "" {
		return nil, errors.New("缺少 Authorization 头")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// 使用 ParseWithClaims 明确指定 MapClaims
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//this.Logger.PrintInfo("jwt secret: ", config.AppCfg.Runtime.JwtSecret)
		return []byte(config.AppCfg.Runtime.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("无效的 Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("无法解析 Token Claims")
	}
	//log.Println(claims)
	return claims, nil
}
