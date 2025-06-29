package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/admin/dto"
	"strconv"
	"strings"
)

// FetchAllMerchants 获取所有商铺信息
func (this *AdminServices) FetchAllMerchants(ctx *gin.Context) {
	var paramsData dto.FetchAllMerchantsRequestDto
	if err := ctx.ShouldBindQuery(&paramsData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "查询参数错误: " + err.Error(),
		})
		return
	}

	this.utils.Logger.PrintInfo("1", paramsData)

	// 默认分页设置
	if paramsData.Page <= 0 {
		paramsData.Page = 1
	}
	if paramsData.Size <= 0 {
		paramsData.Size = 10
	}
	offset := (paramsData.Page - 1) * paramsData.Size

	query := dao.DbDao.Model(&models.Merchant{})

	// 添加搜索条件
	if paramsData.Search != "" {
		like := "%" + paramsData.Search + "%"
		switch paramsData.SearchAs {
		case "name":
			query = query.Where("merchant_name LIKE ?", like)
		case "phone_number":
			query = query.Where("phone_number LIKE ?", like)
		default:
			// 如果 search_as 不合法或未设置，可选行为：忽略或全字段模糊搜索
		}
	}

	// 排序
	sortOrder := "id DESC" // 默认按 id 倒序
	if strings.ToUpper(paramsData.Sort) == "ASC" {
		sortOrder = "id ASC"
	}
	query = query.Order(sortOrder)

	// 查询总数
	var count int64
	if err := query.Count(&count).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询总数失败",
		})
		return
	}

	// 查询数据
	var merchantList []models.Merchant
	if err := query.Offset(int(offset)).Limit(int(paramsData.Size)).Find(&merchantList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询商家失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"merchants": merchantList,
		"count":     count,
		"page":      paramsData.Page,
		"size":      paramsData.Size,
	})
}

// CreateNewMerchant 管理员创建的商铺 POST
func (this *AdminServices) CreateNewMerchant(ctx *gin.Context) {
	var postData dto.CreateNewMerchantRequestDto
	if err := ctx.ShouldBindJSON(&postData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 检查用户是否存在
	var user models.User
	if result := dao.DbDao.Where("id = ? and role = ?", postData.UserId, "trader").First(&user); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "指定的用户不存在或该用户无商家权限",
		})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询用户失败: " + result.Error.Error(),
		})
		return
	}

	if !user.Status {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "该商家账号被禁用",
		})
		return
	}

	// 构造新商户
	newMerchant := models.Merchant{
		UserId:       postData.UserId,
		MerchantId:   uuid.New().String(),
		MerchantName: postData.MerchantName,
		Description:  postData.Description,
		Address:      postData.Address,
		LogoUrl:      postData.LogoUrl,
	}

	// 创建商户
	if result := dao.DbDao.Create(&newMerchant); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建商户失败: " + result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "商户创建成功",
		"merchant": newMerchant,
	})
}

// DeleteMerchant 管理员删除商铺
// 并将删除其附属商品以及订单信息
func (this *AdminServices) DeleteMerchant(ctx *gin.Context) {
	// 1. 获取 URL 中的商户 ID
	idStr := ctx.Param("id")
	merchantID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效的商户 ID"})
		return
	}

	// 2. 查询商户是否存在
	var merchantInfo models.Merchant
	result := dao.DbDao.Where("id = ?", merchantID).First(&merchantInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "商户不存在"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败: " + result.Error.Error()})
		return
	}

	// 3. 删除商户
	if err := dao.DbDao.Delete(&merchantInfo).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败: " + err.Error()})
		return
	}

	// 4. 成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "商户删除成功",
		"id":      merchantID,
	})
}
