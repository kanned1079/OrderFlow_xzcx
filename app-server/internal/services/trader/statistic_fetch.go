package trader

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/trader/dto"
	"strconv"
	"time"
)

func (this *TraderServices) FetchMerchantStatistic(ctx *gin.Context) {
	merchantIdStr := ctx.Param("m_id")
	merchantId, err := strconv.ParseInt(merchantIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "商户ID格式错误"})
		return
	}

	var fetchMerchantStatistic dto.FetchMerchantStatisticResponse

	// 商品数量统计
	if result := dao.DbDao.Model(&models.Goods{}).Unscoped().Where("merchant_id = ?", merchantId).Count(&fetchMerchantStatistic.Goods.Total); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询遇到错误: " + err.Error()})
		return
	}
	if result := dao.DbDao.Model(&models.Goods{}).Where("merchant_id = ?", merchantId).Count(&fetchMerchantStatistic.Goods.OnSale); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询遇到错误: " + err.Error()})
		return
	}
	fetchMerchantStatistic.Goods.Deleted = fetchMerchantStatistic.Goods.Total - fetchMerchantStatistic.Goods.OnSale

	// 过去7天统计（从今天往前推6天）
	var ordersCompleted [7]int64
	var ordersFailed [7]int64
	var income [7]float32

	for i := 0; i < 7; i++ {
		dayStart := time.Now().AddDate(0, 0, -i).Truncate(24 * time.Hour)
		dayEnd := dayStart.Add(24 * time.Hour)

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
		ordersCompleted[6-i] = completedCount // 倒序：0 是最早那天，6 是今天

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
		ordersFailed[6-i] = failedCount

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
		income[6-i] = float32(dailyIncome)
	}

	fetchMerchantStatistic.UserOrders.Completed = ordersCompleted
	fetchMerchantStatistic.UserOrders.Failed = ordersFailed
	fetchMerchantStatistic.Income = income

	ctx.JSON(http.StatusOK, fetchMerchantStatistic)
}
