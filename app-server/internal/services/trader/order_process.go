package trader

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/trader/dto"
	"strings"
	"time"
)

func (TraderServices) GetOrderList(ctx *gin.Context) {
	var req dto.GetOrderListRequestDto
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	offset := (req.Page - 1) * req.Size

	sort := strings.ToUpper(req.Sort)
	if sort != "ASC" {
		sort = "DESC"
	}

	if req.IsActive {
		// 查询活跃订单（状态为待接单或处理中）
		var activeOrders []models.Order
		err := dao.DbDao.
			Where("merchant_id = ? AND (status = ? OR status = ?)", req.MerchantId, "pending_accept", "processing").
			Order("created_at DESC").
			Find(&activeOrders).Error

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"orders": activeOrders,
			"total":  len(activeOrders),
			"page":   1,
			"size":   len(activeOrders),
		})
		return
	}

	// 查询所有订单（分页 + 排序）
	var orders []models.Order
	var total int64

	query := dao.DbDao.Model(&models.Order{}).Where("merchant_id = ?", req.MerchantId)

	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询总数失败: " + err.Error()})
		return
	}

	if err := query.Order("created_at " + sort).
		Offset(offset).
		Limit(req.Size).
		Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"page":   req.Page,
		"size":   req.Size,
	})
}

func (TraderServices) GetOrderById(ctx *gin.Context) {
	// 从 param 获取 order_id
	orderId := ctx.Param("order_id")
	if orderId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "订单号不能为空"})
		return
	}

	var existingOrder models.Order
	if err := dao.DbDao.Where("order_id = ?", orderId).First(&existingOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "未找到该订单"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单失败: " + err.Error()})
		}
		return
	}

	var goodsList []models.OrderItem
	if err := dao.DbDao.Where("order_id = ?", orderId).Find(&goodsList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单商品失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "获取成功",
		"order_id": orderId,
		"order":    existingOrder,
		"goods":    goodsList,
	})
}

// CancelOrderByTrader 商户取消订单
func (TraderServices) CancelOrderByTrader(ctx *gin.Context) {
	var postData dto.CancelOrderByTraderRequestDto
	if err := ctx.ShouldBindJSON(&postData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	var existingOrder models.Order
	result := dao.DbDao.
		Where("order_id = ? AND user_id = ? AND merchant_id = ?", postData.OrderId, postData.UserId, postData.MerchantId).
		First(&existingOrder)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "订单不存在"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查找订单失败: " + result.Error.Error()})
		return
	}

	if existingOrder.Status != "pending_accept" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "该订单已被接单或处理，无法取消"})
		return
	}

	// 启动事务，若有库存扣减，此处可回滚库存
	tx := dao.DbDao.Begin()

	if err := tx.Model(&models.Order{}).
		Where("order_id = ?", postData.OrderId).
		Updates(map[string]interface{}{
			"status":         "cancelled",
			"failure_reason": fmt.Sprintf("%s（取消于 %s）", postData.CancelReason, time.Now().Format("2006-01-02 15:04:05")),
		}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "取消订单失败: " + err.Error()})
		return
	}

	// 查询订单项，准备库存回退
	var orderItems []models.OrderItem
	if err := tx.Where("order_id = ?", existingOrder.OrderId).Find(&orderItems).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单商品项失败: " + err.Error()})
		return
	}

	// 回退库存
	for _, item := range orderItems {
		if err := tx.Model(&models.Goods{}).
			Where("id = ?", item.GoodsId).
			UpdateColumn("residue", gorm.Expr("residue + ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("商品库存回退失败（商品ID: %d）: %s", item.GoodsId, err.Error()),
			})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "提交取消操作失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "订单取消成功",
		"order_id": existingOrder.OrderId,
	})
}

func (TraderServices) AcceptOrderByTrader(ctx *gin.Context) {
	var req dto.AcceptOrderByTraderRequestDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	var existingOrder models.Order
	result := dao.DbDao.Where("order_id = ? AND merchant_id = ?", req.OrderId, req.MerchantId).First(&existingOrder)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "订单不存在或不属于该商户"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单失败: " + result.Error.Error()})
		return
	}

	if existingOrder.Status != "pending_accept" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "当前订单状态不可接单"})
		return
	}

	if err := dao.DbDao.Model(&models.Order{}).
		Where("order_id = ?", req.OrderId).
		Update("status", "processing").Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "接单失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "接单成功",
		"order_id": req.OrderId,
	})
}

func (TraderServices) CompleteOrderByTrader(ctx *gin.Context) {
	var req dto.CompleteOrderByTraderRequestDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	var existingOrder models.Order
	result := dao.DbDao.Where("order_id = ? AND merchant_id = ?", req.OrderId, req.MerchantId).First(&existingOrder)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "订单不存在或不属于该商户"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单失败: " + result.Error.Error()})
		return
	}

	if existingOrder.Status != "processing" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "该订单当前状态不可标记为完成"})
		return
	}

	if err := dao.DbDao.Model(&models.Order{}).
		Where("order_id = ?", req.OrderId).
		Update("status", "completed_unreviewed").Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "订单完成失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "订单完成成功",
		"order_id": req.OrderId,
	})
}
