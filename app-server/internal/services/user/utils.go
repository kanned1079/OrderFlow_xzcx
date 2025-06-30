package user

import (
	"fmt"
	"github.com/google/uuid"
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
