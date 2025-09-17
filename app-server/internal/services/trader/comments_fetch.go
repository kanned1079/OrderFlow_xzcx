package trader

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"strconv"
	"strings"
)

// GetCommentListByMerchantId GET: /api/v1/trader/comments/:m_id?page=&size=&sort=&sort_as=
func (this *TraderServices) GetCommentListByMerchantId(ctx *gin.Context) {
	mIdStr := ctx.Param("m_id")
	merchantId, err := strconv.ParseInt(mIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "地址ID格式错误"})
		return
	}

	// 请求参数
	query := &struct {
		Page   int    `form:"page" json:"page"`
		Size   int    `form:"size" json:"size"`
		Sort   string `form:"sort" json:"sort"`       // ASC / DESC
		SortAs string `form:"sort_as" json:"sort_as"` // 排序字段
	}{}

	if err := ctx.ShouldBindQuery(query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "参数绑定错误: " + err.Error(),
		})
		return
	}

	// 校验排序方式
	sort := strings.ToUpper(query.Sort)
	if sort != "ASC" && sort != "DESC" {
		sort = "DESC"
	}

	// 默认排序字段
	sortField := "id"
	if query.SortAs != "" {
		sortField = query.SortAs
	}

	// 默认分页
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 10
	}
	offset := (query.Page - 1) * query.Size

	// 查询总数
	var total int64
	if err := dao.DbDao.Model(&models.Comment{}).
		Where("merchant_id = ?", merchantId).
		Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询总数失败: " + err.Error(),
		})
		return
	}

	// 查询分页数据
	var list []models.Comment
	if err := dao.DbDao.Model(&models.Comment{}).
		Where("merchant_id = ?", merchantId).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: sortField},
			Desc:   sort == "DESC",
		}).
		Offset(offset).Limit(query.Size).
		Find(&list).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"comments": list,
		"total":    total,
		"page":     query.Page,
		"size":     query.Size,
		"message":  "success",
	})
}
