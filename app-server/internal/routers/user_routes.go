package routers

import (
	"github.com/gin-gonic/gin"
	"stay-server/internal/middlewares"
	"stay-server/internal/services/user"
)

func (this *GatewayApp) RegisterUserRoutes(v1 *gin.RouterGroup) {
	userGroup := v1.Group("/user", middlewares.RequireRole("user"))
	//userGroup := v1.Group("/user")
	var userService user.UserServices

	userGroup.PUT("/:u_id/password/update", userService.UpdateUserPassword)

	userGroup.GET("/merchants", userService.FetchMerchants)
	userGroup.GET("/goods", userService.FetchGoodsListAsCategory)

	userGroup.GET("/order/:u_id", userService.GetUserOrderList)
	userGroup.GET("/order/details/:order_id", userService.GetOrderDetails)
	userGroup.POST("/order", userService.CommitNewOrder)
	userGroup.PUT("/order", userService.CancelOrderByUser)

	userGroup.POST("/comment", userService.CommitCommentByOrderId)
	userGroup.GET("/comment/:m_id", userService.FetchCommentListByMId)
	userGroup.DELETE("/comment/:c_id", userService.DeleteMyComment)

	userGroup.GET("/address/:u_id", userService.FetchAddressLstByUserId) // /api/v1/user/address/:u_id?page=&size=
	userGroup.GET("/address/details/:address_id", userService.FetchAddressById)
	userGroup.POST("/address", userService.AddNewAddress)                 // /api/v1/user/address
	userGroup.PUT("/address/:u_id", userService.UpdateAddressByUserId)    // /api/v1/user/address/:u_id
	userGroup.DELETE("/address/:u_id/:id", userService.DeleteAddressById) // /api/v1/user/address/:u_id/:id

}
