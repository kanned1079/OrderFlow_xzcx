package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/user/dto"
	"time"
)

// GetOrderDetails 获取用户提交的订单细节
func (this *UserServices) GetOrderDetails(ctx *gin.Context) {
	orderId := ctx.Param("order_id")
	if orderId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "缺少订单号"})
		return
	}

	var existingOrder models.Order
	if err := dao.DbDao.Where("order_id = ?", orderId).First(&existingOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "订单不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单失败: " + err.Error()})
		}
		return
	}

	// 查询所有订单项（包括商品快照）
	var orderItems []models.OrderItem
	if err := dao.DbDao.Where("order_id = ?", existingOrder.OrderId).Find(&orderItems).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单商品项失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取订单信息成功",
		"order":   existingOrder,
		"items":   orderItems,
		"count":   len(orderItems),
	})
}

// CommitNewOrder 用户提交新订单
func (this *UserServices) CommitNewOrder(ctx *gin.Context) {
	var postData dto.UserCommitOrderRequestDto
	if err := ctx.ShouldBindJSON(&postData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	if len(postData.GoodsList) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "商品列表不能为空"})
		return
	}

	orderId := this.generateOrderId(postData.UserId, postData.MerchantId)
	newOrder := models.Order{
		OrderId:    orderId,
		UserId:     postData.UserId,
		MerchantId: postData.MerchantId,
		AddressId:  postData.AddressId,
		Status:     "pending_accept", // pending_accept processing completed_unreviewed completed_reviewed cancelled
	}

	tx := dao.DbDao.Begin()

	if err := tx.Create(&newOrder).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "创建订单失败: " + err.Error()})
		return
	}

	for _, item := range postData.GoodsList {
		var goodInfo models.Goods
		this.utils.Logger.PrintInfo("goodsId: ", item.GoodsId, ", MId: ", postData.MerchantId)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}). // 锁住该行
										Where("id = ? AND merchant_id = ?", item.GoodsId, postData.MerchantId).
										First(&goodInfo).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("商品不存在: id=%d", item.GoodsId)})
			return
		}

		if int64(goodInfo.Residue) < item.Count {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("商品 [%s] 库存不足，仅剩 %d 件", goodInfo.GoodsName, goodInfo.Residue),
			})
			return
		}

		// 扣减库存
		if err := tx.Model(&models.Goods{}).
			Where("id = ?", item.GoodsId).
			UpdateColumn("residue", gorm.Expr("residue - ?", item.Count)).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "库存更新失败: " + err.Error()})
			return
		}

		// 创建订单项
		orderItem := models.OrderItem{
			OrderId:    orderId,
			UserId:     postData.UserId,
			MerchantId: postData.MerchantId,
			GoodsId:    item.GoodsId,
			GoodsName:  goodInfo.GoodsName,
			LogoUrl:    goodInfo.LogoUrl,
			Price:      goodInfo.Price,
			Quantity:   int64(item.Count),
			Total:      goodInfo.Price * float32(item.Count),
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "订单项创建失败: " + err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "订单提交失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "新订单提交成功",
		"order_id": orderId,
	})
}

// CancelOrderByUser 用户取消订单
func (UserServices) CancelOrderByUser(ctx *gin.Context) {
	var postData dto.CancelOrderByUserRequestDto
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
			"failure_reason": fmt.Sprintf("%s（取消于 %s）", "用户自主取消", time.Now().Format("2006-01-02 15:04:05")),
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
