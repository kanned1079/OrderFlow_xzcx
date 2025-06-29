package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stay-server/utils"
)

func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	var u utils.Utils
	return func(ctx *gin.Context) {
		claims, err := u.ExtractTokenClaims(ctx)
		//u.Logger.PrintInfo(claims)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未授权: " + err.Error()})
			ctx.Abort()
			return
		}
		//u.Logger.PrintInfo(claims["role"])
		roleClaim, ok := claims["role"].(string)
		if !ok {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Token 中缺少角色信息"})
			ctx.Abort()
			return
		}

		for _, r := range requiredRoles {
			if roleClaim == r {
				ctx.Next()
				return
			}
		}

		ctx.JSON(http.StatusForbidden, gin.H{"message": "无权限访问"})
		ctx.Abort()
	}
}

func RequireAuth() gin.HandlerFunc {
	var u utils.Utils
	return func(ctx *gin.Context) {
		claims, err := u.ExtractTokenClaims(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未授权: " + err.Error()})
			ctx.Abort()
			return
		}

		// ✅ 可选：将 claims 放进上下文，方便后续使用
		ctx.Set("claims", claims)

		ctx.Next()
	}
}
