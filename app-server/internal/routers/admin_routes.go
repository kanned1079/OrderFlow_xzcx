package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"stay-server/internal/middlewares"
	"stay-server/internal/services/admin"
)

func (this *GatewayApp) RegisterAdminRoutes(v1 *gin.RouterGroup) {
	var adminService admin.AdminServices
	var adminGrp = v1.Group("/admin", middlewares.RequireRole("admin"))

	adminGrp.GET("/merchants", adminService.FetchAllMerchants)
	adminGrp.POST("/merchants", adminService.CreateNewMerchant)
	adminGrp.PUT("/merchants/:id", adminService.UpdateMerchantById)
	adminGrp.DELETE("merchants/:id", adminService.DeleteMerchant)

	// 注意这个获取的是商家的列表 而不是商家开的店的列表
	adminGrp.GET("/traders", adminService.FetchAllTraders)
	adminGrp.POST("traders")

	adminGrp.GET("/statistic", adminService.FetchAdminDashboardStatistic)

	adminGrp.GET("/image/clear", func(ctx *gin.Context) {
		path := "./uploads"

		// 读取目录下的所有文件/文件夹
		files, err := os.ReadDir(path)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "读取目录失败", "error": err.Error()})
			return
		}

		// 遍历并删除
		for _, file := range files {
			filePath := filepath.Join(path, file.Name())
			err = os.RemoveAll(filePath) // RemoveAll 可删除文件或文件夹
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败", "file": filePath, "error": err.Error()})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "清理完成"})
	})
}
