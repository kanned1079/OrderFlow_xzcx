package routers

import (
	"github.com/gin-gonic/gin"
	"stay-server/internal/services/user"
)

func (this *GatewayApp) RegisterUserRoutes(v1 *gin.RouterGroup) {
	userGroup := v1.Group("/user")
	var userService user.UserServices

	userGroup.POST("/login", userService.Login)
	userGroup.POST("/register", userService.BindWechatAccountAndLogin)

}
