package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/user/dto"
	"time"
)

func (UserServices) FetchGoodsListAsCategory(ctx *gin.Context) {
	var req dto.FetchGoodsListAsCategoryRequestDto
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	var categories []models.Category
	if err := dao.DbDao.Where("merchant_id = ?", req.MerchantId).Order("created_at DESC").Find(&categories).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询分类失败: " + err.Error()})
		return
	}

	type ResponseGoodsList struct {
		Title     string         `json:"title"`
		CreatedAt time.Time      `json:"created_at"`
		Goods     []models.Goods `json:"goods"`
	}

	var response []ResponseGoodsList

	for _, cat := range categories {
		var goods []models.Goods
		query := dao.DbDao.Where("merchant_id = ? AND category_id = ?", req.MerchantId, cat.Id)

		if req.GoodsName != "" {
			query = query.Where("goods_name LIKE ? or description LIKE ?", "%"+req.GoodsName+"%", "%"+req.GoodsName+"%")
		}

		if err := query.Find(&goods).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询商品失败: " + err.Error()})
			return
		}

		response = append(response, ResponseGoodsList{
			Title:     cat.Title,
			CreatedAt: cat.CreatedAt,
			Goods:     goods,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    response,
	})
}
