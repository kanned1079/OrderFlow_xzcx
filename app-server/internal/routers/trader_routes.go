package routers

import (
	"github.com/gin-gonic/gin"
	"stay-server/internal/services/trader"
)

func (this *GatewayApp) RegisterTraderRoutes(v1 *gin.RouterGroup) {
	var traderService trader.TraderServices
	//traderGrp := v1.Group("/trader", middlewares.RequireRole("trader"))
	traderGrp := v1.Group("/trader")

	traderGrp.GET("/goods/:m_id", traderService.GetGoodsList)
	traderGrp.POST("/goods", traderService.AddNewGoods)
	traderGrp.PUT("/goods/:m_id", traderService.EditGoodsInfo)
	traderGrp.DELETE("/goods/:m_id/:id", traderService.DeleteGoods)

	traderGrp.GET("/category/:m_id", traderService.GetCategoryList)
	traderGrp.POST("/category", traderService.AddNewCategory)
	traderGrp.PUT("/category/:m_id", traderService.EditCategory)
	traderGrp.DELETE("/category/:m_id/:id", traderService.DeleteCategory)

	traderGrp.GET("order/v2/:m_id", traderService.GetOrderList)           // 获取订单列表
	traderGrp.GET("order/:order_id", traderService.GetOrderById)          // 获取订单详情
	traderGrp.PUT("/order/cancel", traderService.CancelOrderByTrader)     // 取消订单
	traderGrp.PUT("/order/accept", traderService.AcceptOrderByTrader)     // 接单
	traderGrp.PUT("/order/complete", traderService.CompleteOrderByTrader) // 将订单完成

	traderGrp.GET("/comments/:m_id", traderService.GetCommentListByMerchantId)

	traderGrp.GET("/address/:a_id", traderService.FetchAddressInfoByAddressId)

	traderGrp.GET("/statistic/:m_id", traderService.FetchMerchantStatistic)

}
