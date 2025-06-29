package trader

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/trader/dto"
	"strconv"
)

// GetCategoryList 获取所有类别列表
func (this *TraderServices) GetCategoryList(ctx *gin.Context) {
	var query dto.GetCategoryListRequestDto
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 10
	}
	offset := (query.Page - 1) * query.Size

	db := dao.DbDao.Model(&models.Category{})
	if query.CategoryTitle != "" {
		db = db.Where("title LIKE ?", "%"+query.CategoryTitle+"%")
	}

	var count int64
	db.Count(&count)

	var list []models.Category
	if err := db.Order("id DESC").Offset(offset).Limit(query.Size).Find(&list).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"categories": list,
		"total":      count,
		"page":       query.Page,
		"size":       query.Size,
	})
}

// AddNewCategory 创建新的分类
func (this *TraderServices) AddNewCategory(ctx *gin.Context) {
	var post dto.AddNewCategoryRequestDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	// 分类名重复检查
	var exists int64
	dao.DbDao.Model(&models.Category{}).
		Where("merchant_id = ? AND title = ?", post.MerchantId, post.CategoryTitle).
		Count(&exists)
	if exists > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"message": "分类名称已存在"})
		return
	}

	newCategory := models.Category{
		MerchantId: post.MerchantId,
		Title:      post.CategoryTitle,
	}

	if err := dao.DbDao.Create(&newCategory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "分类创建成功",
		"category": newCategory,
	})
}

// EditCategory 编辑分类标签
func (this *TraderServices) EditCategory(ctx *gin.Context) {
	var post dto.EditCategoryRequestDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	var category models.Category
	result := dao.DbDao.Where("merchant_id = ? AND title = ?", post.MerchantId, post.CategoryTitle).First(&category)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "分类不存在"})
		return
	}

	// 可以在这里修改字段，例如重命名
	category.Title = post.CategoryTitle

	if err := dao.DbDao.Save(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "分类更新成功"})
}

// DeleteCategory 删除分类
func (this *TraderServices) DeleteCategory(ctx *gin.Context) {
	mIDStr := ctx.Param("m_id")
	categoryIDStr := ctx.Param("id")

	mId, err := strconv.ParseInt(mIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效的商户ID"})
		return
	}
	categoryId, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效的分类ID"})
		return
	}

	var category models.Category
	result := dao.DbDao.Where("id = ? AND merchant_id = ?", categoryId, mId).First(&category)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "分类不存在或无权删除"})
		return
	}

	if err := dao.DbDao.Delete(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "分类删除成功"})
}
