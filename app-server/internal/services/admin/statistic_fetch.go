package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/admin/dto"
	"time"
)

func (this *AdminServices) FetchAdminDashboardStatistic(ctx *gin.Context) {
	var resp dto.FetchAdminDashboardStatisticResponseDto

	now := time.Now()
	yesterdayStart := now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
	todayStart := now.Truncate(24 * time.Hour)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	sevenDaysAgo := todayStart.AddDate(0, 0, -6)

	// 商户统计
	if err := dao.DbDao.Model(&models.Merchant{}).Count(&resp.RegisteredMerchants.Total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询商户总数失败", "error": err.Error()})
		return
	}
	if err := dao.DbDao.Model(&models.Merchant{}).
		Where("created_at >= ? AND created_at < ?", yesterdayStart, todayStart).
		Count(&resp.RegisteredMerchants.Yesterday).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询昨日商户数失败", "error": err.Error()})
		return
	}
	if err := dao.DbDao.Model(&models.Merchant{}).
		Where("created_at >= ?", monthStart).
		Count(&resp.RegisteredMerchants.ThisMonth).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询本月商户数失败", "error": err.Error()})
		return
	}

	// 用户统计
	if err := dao.DbDao.Model(&models.User{}).Count(&resp.RegisteredUsers.Total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询用户总数失败", "error": err.Error()})
		return
	}
	if err := dao.DbDao.Model(&models.User{}).
		Where("created_at >= ? AND created_at < ?", yesterdayStart, todayStart).
		Count(&resp.RegisteredUsers.Yesterday).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询昨日用户数失败", "error": err.Error()})
		return
	}
	if err := dao.DbDao.Model(&models.User{}).
		Where("created_at >= ?", monthStart).
		Count(&resp.RegisteredUsers.ThisMonth).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询本月用户数失败", "error": err.Error()})
		return
	}

	// 最近 7 天订单统计
	var orderStats []struct {
		Date  time.Time
		Count int64
	}
	if err := dao.DbDao.Model(&models.Order{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", sevenDaysAgo).
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&orderStats).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单统计失败", "error": err.Error()})
		return
	}

	orderMap := make(map[string]int64)
	for _, item := range orderStats {
		orderMap[item.Date.Format("2006-01-02")] = item.Count
	}
	for i := 0; i < 7; i++ {
		day := sevenDaysAgo.AddDate(0, 0, i).Format("2006-01-02")
		resp.UserOrdersOverview = append(resp.UserOrdersOverview, orderMap[day])
	}

	// 最近 7 天商户注册统计
	var merchantStats []struct {
		Date  time.Time
		Count int64
	}
	if err := dao.DbDao.Model(&models.Merchant{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", sevenDaysAgo).
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&merchantStats).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询商户注册统计失败", "error": err.Error()})
		return
	}

	merchantMap := make(map[string]int64)
	for _, item := range merchantStats {
		merchantMap[item.Date.Format("2006-01-02")] = item.Count
	}
	for i := 0; i < 7; i++ {
		day := sevenDaysAgo.AddDate(0, 0, i).Format("2006-01-02")
		resp.MerchantsOverview = append(resp.MerchantsOverview, merchantMap[day])
	}

	ctx.JSON(http.StatusOK, &resp)
}
