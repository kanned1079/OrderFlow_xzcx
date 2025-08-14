package routers

import (
	"github.com/gin-gonic/gin"
	"stay-server/internal/middlewares"
	"stay-server/internal/services/admin"
)

func (this *GatewayApp) RegisterAdminRoutes(v1 *gin.RouterGroup) {
	var adminService admin.AdminServices
	var adminGrp = v1.Group("/admin", middlewares.RequireRole("admin"))

	adminGrp.GET("/merchants", adminService.FetchAllMerchants)
	adminGrp.POST("/merchants", adminService.CreateNewMerchant)
	adminGrp.DELETE("merchants/:id", adminService.DeleteMerchant)

	// 注意这个获取的是商家的列表 而不是商家开的店的列表
	adminGrp.GET("/traders", adminService.FetchAllTraders)
	adminGrp.POST("traders")

	adminGrp.GET("/statistic", adminService.FetchAdminDashboardStatistic)
}
