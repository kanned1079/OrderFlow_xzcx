package admin

import (
	"errors"
	"fmt"
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
		query = query.Where("merchant_name LIKE ?", like)
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
//
//	func (this *AdminServices) CreateNewMerchant(ctx *gin.Context) {
//		var postData dto.CreateNewMerchantRequestDto
//		if err := ctx.ShouldBindJSON(&postData); err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{
//				"message": "参数绑定失败: " + err.Error(),
//			})
//			return
//		}
//
//		// 检查用户是否存在
//		var user models.User
//		//if result := dao.DbDao.Where("id = ? and role = ?", postData.UserId, "trader").First(&user); errors.Is(result.Error, gorm.ErrRecordNotFound) {
//		if result := dao.DbDao.Where("phone_number = ?", postData.PhoneNumber).First(&user); errors.Is(result.Error, gorm.ErrRecordNotFound) {
//			ctx.JSON(http.StatusNotFound, gin.H{
//				"message": "指定的用户不存在",
//			})
//			return
//			//if result := dao.DbDao.Model(&models.User{}).Where("phone_number = ?").First(&user); res
//		} else if result.Error != nil {
//			ctx.JSON(http.StatusInternalServerError, gin.H{
//				"message": "查询用户失败: " + result.Error.Error(),
//			})
//			return
//		}
//
//		if !user.Status {
//			ctx.JSON(http.StatusInternalServerError, gin.H{
//				"message": "该商家账号被禁用",
//			})
//			return
//		}
//
//		// 构造新商户
//		newMerchant := models.Merchant{
//			//UserId:       postData.UserId,
//			UserId:       user.Id,
//			MerchantId:   uuid.New().String(),
//			MerchantName: postData.MerchantName,
//			Description:  postData.Description,
//			Address:      postData.Address,
//			LogoUrl:      postData.LogoUrl,
//		}
//
//		// 创建商户
//		if result := dao.DbDao.Create(&newMerchant); result.Error != nil {
//			ctx.JSON(http.StatusInternalServerError, gin.H{
//				"message": "创建商户失败: " + result.Error.Error(),
//			})
//			return
//		}
//
//		user.Role = "trader"
//		if result := dao.DbDao.Model(&models.User{}).Save(user); result.Error != nil {
//			ctx.JSON(http.StatusInternalServerError, gin.H{
//				"message": "修改用户角色失败: " + result.Error.Error(),
//			})
//		}
//
//		ctx.JSON(http.StatusOK, gin.H{
//			"message":  "商户创建成功",
//			"merchant": newMerchant,
//		})
//	}
func (this *AdminServices) CreateNewMerchant(ctx *gin.Context) {
	var postData dto.CreateNewMerchantRequestDto
	if err := ctx.ShouldBindJSON(&postData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	err := dao.DbDao.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否存在
		var user models.User
		if result := tx.Where("phone_number = ?", postData.PhoneNumber).First(&user); errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "指定的用户不存在",
			})
			return fmt.Errorf("user_not_found") // 返回错误触发回滚
		} else if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "查询用户失败: " + result.Error.Error(),
			})
			return result.Error
		}

		if !user.Status {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "该商家账号被禁用",
			})
			return fmt.Errorf("user_disabled")
		}

		// 构造新商户
		newMerchant := models.Merchant{
			UserId:       user.Id,
			MerchantId:   uuid.New().String(),
			MerchantName: postData.MerchantName,
			Description:  postData.Description,
			Address:      postData.Address,
			LogoUrl:      postData.LogoUrl,
		}

		// 创建商户
		if result := tx.Create(&newMerchant); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "创建商户失败: " + result.Error.Error(),
			})
			return result.Error
		}

		// 修改用户角色
		user.Role = "trader"
		if result := tx.Model(&models.User{}).Where("id = ?", user.Id).Update("role", user.Role); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "修改用户角色失败: " + result.Error.Error(),
			})
			return result.Error
		}

		// 返回数据（注意不要在 Transaction 内直接返回 ctx.JSON，避免多次写响应）
		ctx.Set("merchant", newMerchant)
		return nil
	})

	if err != nil {
		// 如果是我们主动返回的业务错误，就不重复输出
		if err.Error() == "user_not_found" || err.Error() == "user_disabled" {
			return
		}
		// 其它错误已在事务内输出过
		return
	}

	// 成功
	if merchant, ok := ctx.Get("merchant"); ok {
		ctx.JSON(http.StatusOK, gin.H{
			"message":  "商户创建成功",
			"merchant": merchant,
		})
	}
}

// UpdateMerchantById 修改商户信息
func (this *AdminServices) UpdateMerchantById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	merchantID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效的商户 ID"})
		return
	}

	// 绑定请求参数
	var postData dto.UpdateMerchantByIdRequestDto
	if err := ctx.ShouldBindJSON(&postData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 开启事务
	tx := dao.DbDao.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "无法开启事务"})
		return
	}

	var merchant models.Merchant
	// 查询商户是否存在
	if err := tx.Where("id = ?", merchantID).First(&merchant).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "商户不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询商户失败"})
		}
		return
	}

	// 更新字段
	merchant.MerchantName = postData.MerchantName
	merchant.Description = postData.Description
	merchant.LogoUrl = postData.LogoUrl
	merchant.Address = postData.Address

	if err := tx.Save(&merchant).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "更新商户失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "提交事务失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    merchant,
	})
}

// DeleteMerchant 管理员删除商铺
// 并将删除其附属商品以及订单信息
//
//	func (this *AdminServices) DeleteMerchant(ctx *gin.Context) {
//		// 1. 获取 URL 中的商户 ID
//		idStr := ctx.Param("id")
//		merchantID, err := strconv.ParseInt(idStr, 10, 64)
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效的商户 ID"})
//			return
//		}
//
//		// 2. 查询商户是否存在
//		var merchantInfo models.Merchant
//		result := dao.DbDao.Where("id = ?", merchantID).First(&merchantInfo)
//		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//			ctx.JSON(http.StatusNotFound, gin.H{"message": "商户不存在"})
//			return
//		} else if result.Error != nil {
//			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败: " + result.Error.Error()})
//			return
//		}
//
//		// 3. 删除商户
//		if err := dao.DbDao.Delete(&merchantInfo).Error; err != nil {
//			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败: " + err.Error()})
//			return
//		}
//
//		dao.DbDao.Model(&models.User{}).Where("id = ?", merchantInfo.UserId).Update("role", "user")
//
//		// 4. 成功响应
//		ctx.JSON(http.StatusOK, gin.H{
//			"message": "商户删除成功",
//			"id":      merchantID,
//		})
//	}
func (this *AdminServices) DeleteMerchant(ctx *gin.Context) {
	// 获取商户 ID
	idStr := ctx.Param("id")
	merchantID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "无效的商户 ID"})
		return
	}

	err = dao.DbDao.Transaction(func(tx *gorm.DB) error {
		// 查询商户是否存在
		var merchantInfo models.Merchant
		if result := tx.Where("id = ?", merchantID).First(&merchantInfo); errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "商户不存在"})
			return fmt.Errorf("merchant_not_found")
		} else if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败: " + result.Error.Error()})
			return result.Error
		}

		// 删除商户
		if err := tx.Delete(&merchantInfo).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败: " + err.Error()})
			return err
		}

		// 还原用户角色为普通用户
		if err := tx.Model(&models.User{}).Where("id = ?", merchantInfo.UserId).Update("role", "user").Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "用户角色还原失败: " + err.Error()})
			return err
		}

		// 保存删除成功的 ID
		ctx.Set("deleted_id", merchantID)
		return nil
	})

	if err != nil {
		if err.Error() == "merchant_not_found" {
			return // 业务错误已响应
		}
		return // 其它错误事务内已响应
	}

	// 成功响应
	if deletedID, ok := ctx.Get("deleted_id"); ok {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "商户删除成功",
			"id":      deletedID,
		})
	}
}
