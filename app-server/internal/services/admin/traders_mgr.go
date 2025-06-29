package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/admin/dto"
)

func (this *AdminServices) FetchAllTraders(ctx *gin.Context) {
	var paramsData dto.FetchAllTradersRequestDto
	if err := ctx.ShouldBindQuery(&paramsData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "查询参数错误: " + err.Error(),
		})
		return
	}

	if paramsData.Page <= 0 {
		paramsData.Page = 1
	}
	if paramsData.Size <= 0 {
		paramsData.Size = 10
	}
	offset := (paramsData.Page - 1) * paramsData.Size

	query := dao.DbDao.Model(&models.User{}).Where("role = ?", "trader")

	// 可选搜索：手机号
	if paramsData.PhoneNumber != "" {
		query = query.Where("phone_number LIKE ?", "%"+paramsData.PhoneNumber+"%")
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询总数失败"})
		return
	}

	var traderList []models.User
	if err := query.Offset(offset).Limit(paramsData.Size).Find(&traderList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}

	// 删除 password 字段
	for i := range traderList {
		traderList[i].Password = ""
	}

	ctx.JSON(http.StatusOK, gin.H{
		"traders": traderList,
		"count":   count,
		"page":    paramsData.Page,
		"size":    paramsData.Size,
	})
}
