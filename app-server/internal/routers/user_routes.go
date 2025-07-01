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

	userGroup.GET("goods", userService.FetchGoodsListAsCategory)

	userGroup.GET("order/:order_id", userService.GetOrderDetails)
	userGroup.POST("order", userService.CommitNewOrder)
	userGroup.PUT("order", userService.CancelOrderByUser)

	userGroup.POST("comment", userService.CommitCommentByOrderId)
	userGroup.GET("comment/:m_id", userService.FetchCommentListByMId)
	userGroup.DELETE("comment/:c_id", userService.DeleteMyComment)

}
