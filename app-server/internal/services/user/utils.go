package user

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"stay-server/internal/models"
	"stay-server/utils"
	"strings"
	"time"
)

// generateOrderId 生成订单号
func (UserServices) generateOrderId(userId, merchantId int64) string {
	now := time.Now()

	// 日期：250628（YYMMDD）
	datePart := now.Format("060102")

	// 今日 0 点
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	secondsSinceMidnight := int64(now.Sub(todayStart).Seconds()) // 如 54065 表示15:01:05

	// 短UUID（无横线前8位）
	shortUUID := strings.ReplaceAll(uuid.New().String(), "-", "")[:8]

	// 拼接格式：短UUID + 日期 + 秒数 + 用户ID + 商户ID
	orderID := fmt.Sprintf("%s%s%d%d%d", shortUUID, datePart, secondsSinceMidnight, userId, merchantId)

	return orderID
}

// updateMerchantStars 提交评论时更新商家的评分
func (UserServices) updateMerchantStars(tx *gorm.DB, merchantId int64, newStars int8) error {
	utils.Logger{}.PrintInfo(merchantId, newStars)

	var merchant models.Merchant
	// 1. 加锁获取
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", merchantId).
		First(&merchant).Error; err != nil {
		return err
	}

	// 2. 本地计算新平均分
	newCount := merchant.StarsCount + 1
	newAvg := (merchant.AvgStarsFloat*float64(merchant.StarsCount) + float64(newStars)) / float64(newCount)

	// 3. 更新字段
	return tx.Model(&models.Merchant{}).
		Where("id = ?", merchantId).
		Updates(map[string]interface{}{
			"avg_stars":   newAvg,
			"stars_count": newCount,
		}).Error
}

// decreaseMerchantStars 用户删除评论后重新计算商家评分
func (UserServices) decreaseMerchantStars(tx *gorm.DB, merchantId int64, stars int8) error {
	// 1. 加锁查询商家数据
	var merchant models.Merchant
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", merchantId).
		First(&merchant).Error; err != nil {
		return err
	}

	// 2. 计算新评分
	var newAvg float64
	var newCount int64

	if merchant.StarsCount <= 1 {
		newAvg = 0.0
		newCount = 0
	} else {
		newCount = merchant.StarsCount - 1
		newAvg = (merchant.AvgStarsFloat*float64(merchant.StarsCount) - float64(stars)) / float64(newCount)
	}

	// 3. 更新到数据库
	return tx.Model(&models.Merchant{}).
		Where("id = ?", merchantId).
		Updates(map[string]interface{}{
			"avg_stars":   newAvg,
			"stars_count": newCount,
		}).Error
}
