package routers

import (
	"github.com/gin-gonic/gin"
	"stay-server/internal/middlewares"
	"stay-server/internal/services/trader"
)

func (this *GatewayApp) RegisterTraderRoutes(v1 *gin.RouterGroup) {
	var traderService trader.TraderServices
	traderGrp := v1.Group("/trader", middlewares.RequireRole("trader"))

	traderGrp.GET("/goods", traderService.GetGoodsList)
	traderGrp.POST("/goods", traderService.AddNewGoods)
	traderGrp.PUT("/goods", traderService.EditGoodsInfo)
	traderGrp.DELETE("/goods/:m_id/:id", traderService.DeleteGoods)

	traderGrp.GET("/category", traderService.GetCategoryList)
	traderGrp.POST("/category", traderService.AddNewCategory)
	traderGrp.PUT("/category", traderService.EditCategory)
	traderGrp.DELETE("/category/:m_id/:id", traderService.DeleteCategory)

	traderGrp.GET("order", traderService.GetOrderList)
	traderGrp.GET("order/:order_id", traderService.GetOrderById)
	traderGrp.PUT("/order/cancel", traderService.CancelOrderByTrader)
	traderGrp.PUT("/order/accept", traderService.AcceptOrderByTrader)
	traderGrp.PUT("/order/complete", traderService.CompleteOrderByTrader)

	traderGrp.GET("/statistic/:m_id", traderService.FetchMerchantStatistic)

}
