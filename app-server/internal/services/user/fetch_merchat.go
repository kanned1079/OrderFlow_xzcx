package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/user/dto"
	"strings"
)

// FetchMerchants GET: /api/v1/user/merchants?search=&page=&size=&sort_as=&sort=
func (this *UserServices) FetchMerchants(ctx *gin.Context) {
	var searchReq dto.FetchMerchantsRequestDto

	// 绑定 query 参数
	if err := ctx.ShouldBindQuery(&searchReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "查询参数错误: " + err.Error(),
		})
		return
	}

	// 默认分页
	if searchReq.Page <= 0 {
		searchReq.Page = 1
	}
	if searchReq.Size <= 0 {
		searchReq.Size = 10
	}
	offset := (searchReq.Page - 1) * searchReq.Size

	// 默认排序
	sort := strings.ToUpper(searchReq.Sort)
	if sort != "ASC" && sort != "DESC" {
		sort = "ASC"
	}

	sortField := "id"
	if searchReq.SortAs == "created_at" {
		sortField = "created_at"
	}

	// 构建查询
	query := dao.DbDao.Model(&models.Merchant{})

	if searchReq.Search != "" {
		like := "%" + searchReq.Search + "%"
		query = query.Where("merchant_name LIKE ? OR description LIKE ?", like, like)
	}

	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询商家总数失败",
			"error":   err.Error(),
		})
		return
	}

	// 查询分页数据
	var merchantList []models.Merchant
	if err := query.
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: sortField},
			Desc:   sort == "DESC",
		}).
		Offset(offset).
		Limit(searchReq.Size).
		Find(&merchantList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询商家失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"merchants": merchantList,
		"total":     total,
		"page":      searchReq.Page,
		"size":      searchReq.Size,
		"search":    searchReq.Search,
		"sort":      sort,
		"sort_as":   sortField,
		"message":   "success",
	})
}
