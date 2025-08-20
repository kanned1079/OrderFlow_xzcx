package trader

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/trader/dto"
	"strconv"
	"time"
)

// 生成随机日期
func (this *TraderServices) randomDates(days int) []string {
	now := time.Now()
	dates := make([]string, days)
	for i := 0; i < days; i++ {
		dates[i] = now.AddDate(0, 0, -i).Format("01/02") // MM/DD 格式
	}
	return dates
}

func (this *TraderServices) randStatisticDataGenerate() dto.FetchMerchantStatisticResponse {
	// rand.Seed(time.Now().UnixNano())

	var resp dto.FetchMerchantStatisticResponse

	// 商品计数
	resp.GoodsCount.Total = int64(rand.Intn(1000) + 500)
	resp.GoodsCount.Deleted = int64(rand.Intn(200))
	resp.GoodsCount.OnSale = resp.GoodsCount.Total - resp.GoodsCount.Deleted

	// 用户订单
	resp.UserOrders.TodayCount = int64(rand.Intn(50))
	resp.UserOrders.TodayProcessing = int64(rand.Intn(int(resp.UserOrders.TodayCount)))
	resp.UserOrders.WeekFinishedRate = rand.Float32()

	// 最近 7 天数据
	dates := this.randomDates(7)
	for _, d := range dates {
		resp.UserOrders.CompletedGraph = append(resp.UserOrders.CompletedGraph, dto.OverviewItem{
			Date:  d,
			Value: float64(rand.Intn(50)),
		})
		resp.UserOrders.FailedGraph = append(resp.UserOrders.FailedGraph, dto.OverviewItem{
			Date:  d,
			Value: float64(rand.Intn(10)),
		})
		resp.IncomeGraph = append(resp.IncomeGraph, dto.OverviewItem{
			Date:  d,
			Value: float64(rand.Intn(1000) + 100),
		})
	}

	// 模拟正在处理的订单
	for i := 0; i < rand.Intn(50)+1; i++ {
		setupTime := time.Now().Add(-time.Duration(rand.Intn(48)) * time.Hour)

		resp.ActiveOrderList = append(resp.ActiveOrderList, models.Order{
			OrderId:     fmt.Sprintf("%s%v%d", time.Now().Format("20060102"), time.Now().Unix(), i),
			UserId:      rand.Int63n(1000),
			TotalAmount: float32(rand.Intn(500) + 50),
			MerchantId:  -1,
			AddressId:   -1,
			Status:      []string{"pending", "processing"}[rand.Intn(2)],
			CreatedAt:   setupTime,
			UpdatedAt:   setupTime,
		})
	}

	resp.Message = "mock data (TEST ONLY)"

	return resp
}

// FetchMerchantStatistic 获取商家的统计数据
func (this *TraderServices) FetchMerchantStatistic(ctx *gin.Context) {
	merchantIdStr := ctx.Param("m_id")
	if merchantIdStr == "test" {
		ctx.JSON(http.StatusOK, this.randStatisticDataGenerate())
		return
	}
	merchantId, err := strconv.ParseInt(merchantIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "商户ID格式错误"})
		return
	}
	var fetchMerchantStatistic dto.FetchMerchantStatisticResponse

	// 活跃订单（待接单 + 正在处理）
	if result := dao.DbDao.Model(&models.Order{}).
		Where("merchant_id = ? AND status IN ?", merchantId, []string{"pending_accept", "processing"}).
		Find(&fetchMerchantStatistic.ActiveOrderList); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询活跃订单出错"})
		return
	}

	// 商品数量统计
	if result := dao.DbDao.Model(&models.Goods{}).Unscoped().Where("merchant_id = ?", merchantId).
		Count(&fetchMerchantStatistic.GoodsCount.Total); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询商品总数出错"})
		return
	}
	if result := dao.DbDao.Model(&models.Goods{}).Where("merchant_id = ?", merchantId).
		Count(&fetchMerchantStatistic.GoodsCount.OnSale); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询在售商品出错"})
		return
	}
	fetchMerchantStatistic.GoodsCount.Deleted = fetchMerchantStatistic.GoodsCount.Total - fetchMerchantStatistic.GoodsCount.OnSale

	// -------------------------------
	// 今日订单统计
	// -------------------------------
	todayStart := time.Now().Truncate(24 * time.Hour)
	todayEnd := todayStart.Add(24 * time.Hour)

	// 今天总订单数
	if err := dao.DbDao.Model(&models.Order{}).
		Where("merchant_id = ? AND created_at >= ? AND created_at < ?", merchantId, todayStart, todayEnd).
		Count(&fetchMerchantStatistic.UserOrders.TodayCount).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询今日订单失败"})
		return
	}

	// 今天正在处理中的订单数
	if err := dao.DbDao.Model(&models.Order{}).
		Where("merchant_id = ? AND created_at >= ? AND created_at < ? AND status = ?", merchantId, todayStart, todayEnd, "processing").
		Count(&fetchMerchantStatistic.UserOrders.TodayProcessing).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询今日处理中订单失败"})
		return
	}

	// -------------------------------
	// 最近 7 天统计
	// -------------------------------
	var completedList []dto.OverviewItem
	var failedList []dto.OverviewItem
	var incomeList []dto.OverviewItem

	var totalWeekOrders int64
	var totalWeekCompleted int64

	for i := 0; i < 7; i++ {
		dayStart := time.Now().AddDate(0, 0, -i).Truncate(24 * time.Hour)
		dayEnd := dayStart.Add(24 * time.Hour)
		dayStr := dayStart.Format("01/02") // 格式化为 MM/DD

		// 成功订单数量
		var completedCount int64
		if err := dao.DbDao.Model(&models.Order{}).
			Unscoped().
			Where("merchant_id = ? AND status IN ? AND created_at >= ? AND created_at < ?", merchantId,
				[]string{"completed_unreviewed", "completed_reviewed"}, dayStart, dayEnd).
			Count(&completedCount).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询成功订单失败"})
			return
		}
		completedList = append([]dto.OverviewItem{{Date: dayStr, Value: float64(completedCount)}}, completedList...)
		totalWeekCompleted += completedCount

		// 失败订单数量
		var failedCount int64
		if err := dao.DbDao.Model(&models.Order{}).
			Unscoped().
			Where("merchant_id = ? AND status = ? AND created_at >= ? AND created_at < ?", merchantId,
				"cancelled", dayStart, dayEnd).
			Count(&failedCount).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败订单失败"})
			return
		}
		failedList = append([]dto.OverviewItem{{Date: dayStr, Value: float64(failedCount)}}, failedList...)

		// 每日收入统计
		var dailyIncome float64
		if err := dao.DbDao.Model(&models.Order{}).
			Select("COALESCE(SUM(total_amount), 0)").
			Where("merchant_id = ? AND status IN ? AND created_at >= ? AND created_at < ?", merchantId,
				[]string{"completed_unreviewed", "completed_reviewed"}, dayStart, dayEnd).
			Scan(&dailyIncome).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询收入失败"})
			return
		}
		incomeList = append([]dto.OverviewItem{{Date: dayStr, Value: dailyIncome}}, incomeList...)

		// 周统计总订单数
		var dayOrders int64
		if err := dao.DbDao.Model(&models.Order{}).
			Where("merchant_id = ? AND created_at >= ? AND created_at < ?", merchantId, dayStart, dayEnd).
			Count(&dayOrders).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询周订单数失败"})
			return
		}
		totalWeekOrders += dayOrders
	}

	// 填充结果
	fetchMerchantStatistic.UserOrders.CompletedGraph = completedList
	fetchMerchantStatistic.UserOrders.FailedGraph = failedList
	fetchMerchantStatistic.IncomeGraph = incomeList

	// 周完成率 (完成数 / 总数)
	if totalWeekOrders > 0 {
		fetchMerchantStatistic.UserOrders.WeekFinishedRate = float32(totalWeekCompleted) / float32(totalWeekOrders)
	} else {
		fetchMerchantStatistic.UserOrders.WeekFinishedRate = 0
	}

	fetchMerchantStatistic.Message = "success"

	ctx.JSON(http.StatusOK, fetchMerchantStatistic)
}
