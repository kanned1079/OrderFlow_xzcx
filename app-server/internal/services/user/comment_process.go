package user

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/user/dto"
	"strconv"
	"strings"
	"time"
)

// CommitCommentByOrderId 上传评论
func (this *UserServices) CommitCommentByOrderId(ctx *gin.Context) {
	var commitCommentRequestDto dto.CommitCommentByOrderIdRequestDto

	// 1. 参数绑定校验
	if err := ctx.ShouldBindJSON(&commitCommentRequestDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	if strings.Trim(commitCommentRequestDto.CommentText, " ") == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "评论内容不能为空"})
		return
	}

	if commitCommentRequestDto.Stars <= 0 || commitCommentRequestDto.Stars > 5 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "非法评分"})
		return
	}

	// 2. 开始事务
	tx := dao.DbDao.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "数据库事务开启失败"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		}
	}()

	// 3. 查询订单
	var existingOrder models.Order
	if err := tx.
		Where("order_id = ?", commitCommentRequestDto.OrderId).
		First(&existingOrder).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"message": "订单不存在"})
		return
	}

	// 4. 校验订单状态
	if existingOrder.Status != "completed_unreviewed" {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "该订单暂未完成或已经评价"})
		return
	}

	// 5. 图片转 JSON
	jsonImageUrlList, err := json.Marshal(commitCommentRequestDto.ImagesUrls)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "图片格式化错误: " + err.Error()})
		return
	}

	// 6. 创建评论记录
	newComment := models.Comment{
		OrderId:     existingOrder.OrderId,
		UserId:      commitCommentRequestDto.UserId,
		MerchantId:  existingOrder.MerchantId,
		Stars:       commitCommentRequestDto.Stars,
		CommentText: commitCommentRequestDto.CommentText,
		ImagesUrls:  jsonImageUrlList,
	}
	if err := tx.Create(&newComment).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "评论保存失败: " + err.Error()})
		return
	}

	// 6.5. 更新商家评分
	if err := this.updateMerchantStars(tx, existingOrder.MerchantId, commitCommentRequestDto.Stars); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "评分更新失败: " + err.Error()})
		return
	}

	// 7. 更新订单状态
	existingOrder.Status = "completed_reviewed"
	if err := tx.Save(&existingOrder).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "订单状态更新失败: " + err.Error()})
		return
	}

	// 8. 提交事务
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "事务提交失败: " + err.Error()})
		return
	}

	// 9. 成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "发表评论成功",
		"comment": newComment,
	})
}

// FetchCommentListByMId 通过商户的Id来查询评论列表
func (UserServices) FetchCommentListByMId(ctx *gin.Context) {
	mIdParam := ctx.Param("c_id")
	merchantId, err := strconv.ParseInt(mIdParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "商户 ID 非法"})
		return
	}

	// 2. 从查询参数中绑定分页参数
	var fetchCommentListByMIdRequestDto dto.FetchCommentListByMIdRequestDto
	if err := ctx.ShouldBindQuery(&fetchCommentListByMIdRequestDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "分页参数错误: " + err.Error()})
		return
	}
	if fetchCommentListByMIdRequestDto.Page <= 0 {
		fetchCommentListByMIdRequestDto.Page = 1
	}
	if fetchCommentListByMIdRequestDto.Size <= 0 {
		fetchCommentListByMIdRequestDto.Size = 10
	}
	offset := (fetchCommentListByMIdRequestDto.Page - 1) * fetchCommentListByMIdRequestDto.Size

	// 3. 查询评论总数
	var total int64
	if err := dao.DbDao.Model(&models.Comment{}).
		Where("merchant_id = ?", merchantId).
		Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询总数失败: " + err.Error()})
		return
	}

	// 4. 查询评论分页列表
	var commentList []models.Comment
	if err := dao.DbDao.Model(&models.Comment{}).
		Where("merchant_id = ?", merchantId).
		Order("created_at DESC").
		Limit(fetchCommentListByMIdRequestDto.Size).
		Offset(offset).
		Find(&commentList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询评论失败: " + err.Error()})
		return
	}

	// 5. 成功返回
	ctx.JSON(http.StatusOK, gin.H{
		"comments": commentList,
		"total":    total,
		"page":     fetchCommentListByMIdRequestDto.Page,
		"size":     fetchCommentListByMIdRequestDto.Size,
	})
}

// FetchCommentbyId 通过评论的Id来获取评论细节
func (this *UserServices) FetchCommentbyId(ctx *gin.Context) {
	commentIdParam := ctx.Param("c_id")
	commentId, err := strconv.ParseInt(commentIdParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "非法评论ID"})
		return
	}

	var existingComment models.Comment

	if result := dao.DbDao.Model(&models.Comment{}).Where("id = ?", commentId).First(&existingComment); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "无法找到此id的评论"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"comment": existingComment,
	})
}

//func (this *UserServices) DeleteMyComment(ctx *gin.Context) {
//	commentIdParam := ctx.Param("c_id")
//	commentId, err := strconv.ParseInt(commentIdParam, 10, 64)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"message": "非法评论ID"})
//		return
//	}
//
//	var existingComment models.Comment
//	if result := dao.DbDao.Model(&models.Comment{}).Where("id = ?", commentId).First(&existingComment); errors.Is(result.Error, gorm.ErrRecordNotFound) {
//		ctx.JSON(http.StatusNotFound, gin.H{"message": "没有找到指定的评论，该评论是否不存在或已被删除。"})
//		return
//	} else if result.Error != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询评论失败: " + err.Error()})
//		return
//	}
//
//	if time.Now().Sub(existingComment.CreatedAt) > time.Hour*3 {
//		ctx.JSON(http.StatusRequestTimeout, gin.H{
//			"message": "已超出删除评论的有效时间（3h）,不可操作。",
//		})
//		return
//	}
//
//	if delResult := dao.DbDao.Model(&models.Comment{}).Delete(&existingComment); delResult.RowsAffected == 0 && delResult.Error != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "删除评论失败: " + delResult.Error.Error()})
//		return
//	}
//
//	// 重新计算商家评分
//	//this.decreaseMerchantStars()
//
//	ctx.JSON(http.StatusOK, gin.H{
//		"message":    "评论已经删除",
//		"comment_id": existingComment.Id,
//	})
//
//}

func (this *UserServices) DeleteMyComment(ctx *gin.Context) {
	commentIdParam := ctx.Param("c_id")
	commentId, err := strconv.ParseInt(commentIdParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "非法评论ID"})
		return
	}

	// 开启事务
	tx := dao.DbDao.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "数据库事务开启失败"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		}
	}()

	var existingComment models.Comment
	if result := tx.Where("id = ?", commentId).First(&existingComment); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"message": "没有找到指定的评论，该评论是否不存在或已被删除。"})
		return
	} else if result.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询评论失败: " + result.Error.Error()})
		return
	}

	// 限时删除
	if time.Since(existingComment.CreatedAt) > 3*time.Hour {
		tx.Rollback()
		ctx.JSON(http.StatusRequestTimeout, gin.H{
			"message": "已超出删除评论的有效时间（3h）,不可操作。",
		})
		return
	}

	// 删除评论
	if delResult := tx.Delete(&existingComment); delResult.RowsAffected == 0 || delResult.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "删除评论失败: " + delResult.Error.Error()})
		return
	}

	// 更新商家评分
	if err := this.decreaseMerchantStars(tx, existingComment.MerchantId, existingComment.Stars); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "更新评分失败: " + err.Error()})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "事务提交失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "评论已经删除，并更新评分",
		"comment_id": existingComment.Id,
	})
}
