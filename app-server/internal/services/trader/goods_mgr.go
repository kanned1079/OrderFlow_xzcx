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

// GetGoodsList 商家获取商品列表
func (TraderServices) GetGoodsList(ctx *gin.Context) {
	var paramsData dto.GetGoodsListRequestDto
	if err := ctx.ShouldBindQuery(&paramsData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "查询参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if paramsData.Page <= 0 {
		paramsData.Page = 1
	}
	if paramsData.Size <= 0 {
		paramsData.Size = 10
	}
	if paramsData.Sort != "ASC" {
		paramsData.Sort = "DESC"
	}
	//offset := (paramsData.Page - 1) * paramsData.Size

	query := dao.DbDao.Model(&models.Goods{}).Where("merchant_id = ?", paramsData.MerchantId)

	// 模糊查询商品名
	if paramsData.GoodsName != "" {
		query = query.Where("description LIKE ? OR goods_name LIKE ?", "%"+paramsData.GoodsName+"%", "%"+paramsData.GoodsName+"%")
	}

	// 查询总数
	var count int64
	if err := query.Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询总数失败"})
		return
	}

	// 查询分页数据
	var goodsList []models.Goods
	if err := query.Order("id " + paramsData.Sort).
		Offset((paramsData.Page - 1) * paramsData.Size).
		Limit(paramsData.Size).
		Find(&goodsList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询商品失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"goods": goodsList,
		"total": count,
		"page":  paramsData.Page,
		"size":  paramsData.Size,
	})
}

// AddNewGoods 商户添加商品 但是需要先在分类中设置好
func (this *TraderServices) AddNewGoods(ctx *gin.Context) {
	var postData dto.AddNewGoodsRequestDto
	if err := ctx.ShouldBindJSON(&postData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	if postData.Price <= 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "商品价格不正确",
		})
		return
	}

	//postData.MerchantId // 商户id
	//// 查询该商户信息
	//var existingMerchant models.Merchant
	//
	//existingMerchant.UserId // 查询该用户id
	//

	// 检查商品是否已存在（同一商户下不能重名）
	var existing models.Goods
	result := dao.DbDao.Where("merchant_id = ? AND goods_name = ?", postData.MerchantId, postData.GoodsName).First(&existing)
	if result.Error == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "该商品已存在，请勿重复添加",
		})
		return
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询失败: " + result.Error.Error(),
		})
		return
	}

	// 查询category表 如果存在才可以添加 否责500 提示错误
	exists, err := this.categoryExistsForMerchant(postData.CategoryId, postData.MerchantId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "分类验证失败: " + err.Error()})
		return
	}
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "该分类不存在或不属于该商户"})
		return
	}

	// 构建新商品
	newGoods := models.Goods{
		MerchantId:  postData.MerchantId,
		CategoryId:  postData.CategoryId,
		GoodsName:   postData.GoodsName,
		Description: postData.Description,
		LogoUrl:     postData.LogoUrl,
		Price:       postData.Price,
		Residue:     postData.Residue,
	}

	if err := dao.DbDao.Create(&newGoods).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建商品失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "商品创建成功",
		"goods":   newGoods,
	})
}

// EditGoodsInfo 商家编辑自家商品的信息
func (this *TraderServices) EditGoodsInfo(ctx *gin.Context) {
	var postData dto.EditGoodsRequestDto
	if err := ctx.ShouldBindJSON(&postData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	if postData.Price <= 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "商品价格不正确",
		})
		return
	}

	var goods models.Goods
	result := dao.DbDao.
		Where("id = ? AND merchant_id = ?", postData.Id, postData.MerchantId).
		First(&goods)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "商品不存在或无权操作",
		})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询失败: " + result.Error.Error(),
		})
		return
	}

	exists, err := this.categoryExistsForMerchant(postData.CategoryId, postData.MerchantId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "分类验证失败: " + err.Error()})
		return
	}
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "该分类不存在或不属于该商户"})
		return
	}

	// 执行更新
	update := map[string]interface{}{
		"goods_name":  postData.GoodsName,
		"category_id": postData.CategoryId,
		"description": postData.Description,
		"logo_url":    postData.LogoUrl,
		"price":       postData.Price,
		"residue":     postData.Residue,
	}
	if err := dao.DbDao.Model(&goods).Updates(update).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "商品信息更新成功",
		"goods":   goods,
	})
}

func (TraderServices) DeleteGoods(ctx *gin.Context) {
	// 获取路由参数
	mIDStr := ctx.Param("m_id")
	goodsIDStr := ctx.Param("id")

	// 转换为 int64
	mID, err := strconv.ParseInt(mIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效的商户ID"})
		return
	}
	goodsID, err := strconv.ParseInt(goodsIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效的商品ID"})
		return
	}

	// 查询是否存在该商品，并且属于该商户
	var goods models.Goods
	result := dao.DbDao.Where("id = ? AND merchant_id = ?", goodsID, mID).First(&goods)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "商品不存在或无权限删除"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败: " + result.Error.Error()})
		return
	}

	// 执行删除（软删除）
	if err := dao.DbDao.Delete(&goods).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "商品删除成功",
		"goods_id": goods.Id,
	})
}
