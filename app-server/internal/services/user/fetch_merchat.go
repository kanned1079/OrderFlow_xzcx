package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/user/dto"
)

func (this *UserServices) FetchMerchants(ctx *gin.Context) {
	var searchReq dto.FetchMerchantsRequestDto
	// 绑定 query 参数
	if err := ctx.ShouldBindQuery(&searchReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "查询参数错误: " + err.Error(),
		})
		return
	}

	//this.utils.Logger.PrintInfo("搜索商家: ", searchReq.Search)

	var merchantList []models.Merchant
	query := dao.DbDao.Model(&models.Merchant{})

	if searchReq.Search != "" {
		like := "%" + searchReq.Search + "%"
		query = query.Where("merchant_name LIKE ? OR description LIKE ?", like, like)
	}

	if err := query.Find(&merchantList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询商家失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"merchants": merchantList,
		"count":     len(merchantList),
		"search":    searchReq.Search,
	})
}
