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
		day := sevenDaysAgo.AddDate(0, 0, i)
		dateKey := day.Format("2006-01-02")
		resp.UserOrdersOverview = append(resp.UserOrdersOverview, dto.OverviewItem{
			Date:  day.Format("01/02"), // MM/DD
			Value: orderMap[dateKey],
		})
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
		day := sevenDaysAgo.AddDate(0, 0, i)
		dateKey := day.Format("2006-01-02")
		resp.MerchantsOverview = append(resp.MerchantsOverview, dto.OverviewItem{
			Date:  day.Format("01/02"), // MM/DD
			Value: merchantMap[dateKey],
		})
	}

	resp.P1 = `“绿水青山既是自然财富，又是经济财富。”习近平总书记的这一重要论述，深刻揭示了绿水青山与金山银山的辩证统一关系。人不负青山，青山定不负人。党的十八大以来，以习近平同志为核心的党中央以前所未有的力度抓生态文明建设，一幅幅生态美、产业兴、百姓富的壮美画卷已铺展开来。`

	ctx.JSON(http.StatusOK, &resp)
}
