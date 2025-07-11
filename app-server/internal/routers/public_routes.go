package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stay-server/internal/services/user"
)

func (this *GatewayApp) RegisterPublicRoutes(v1 *gin.RouterGroup) {
	public := v1.Group("/public")
	var userServices user.UserServices

	public.GET("/test", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "dfgbnrdfbs",
		})
	})

	public.POST("/user/login", userServices.Login)
	public.POST("/user/register", userServices.Register)
}
