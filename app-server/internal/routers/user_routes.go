package routers

import (
	"github.com/gin-gonic/gin"
	"stay-server/internal/middlewares"
	"stay-server/internal/services/user"
)

func (this *GatewayApp) RegisterUserRoutes(v1 *gin.RouterGroup) {
	userGroup := v1.Group("/user", middlewares.RequireRole("user"))
	var userService user.UserServices

	userGroup.GET("merchants", userService.FetchMerchants)

}
