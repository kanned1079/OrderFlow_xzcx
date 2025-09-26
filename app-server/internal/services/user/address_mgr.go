package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/user/dto"
	"strconv"
)

// FetchAddressLstByUserId 用户获取地址列表 GET: /api/v1/user/address/:u_id?page=&size=
func (UserServices) FetchAddressLstByUserId(ctx *gin.Context) {
	userIdStr := ctx.Param("u_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}

	// 解析分页参数
	query := struct {
		Page int `form:"page" json:"page"`
		Size int `form:"size" json:"size"`
	}{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var (
		addrList []models.Address
		total    int64
	)

	// 先统计总数
	if err := dao.DbDao.Model(&models.Address{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	db := dao.DbDao.Model(&models.Address{}).Where("user_id = ?", userId)

	// page = -1 => 查询全部
	if query.Page == -1 {
		if err := db.Find(&addrList).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"list":    addrList,
			"message": "success",
			"total":   total,
			"page":    query.Page,
			"size":    total, // 这里 size 可以直接返回 total，表示全部
		})
		return
	}

	// 默认分页参数
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 10
	}
	offset := (query.Page - 1) * query.Size

	// 查询分页数据
	if err := db.Limit(query.Size).Offset(offset).Find(&addrList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"list":    addrList,
		"message": "success",
		"total":   total,
		"page":    query.Page,
		"size":    query.Size,
	})
}

func (this *UserServices) FetchAddressById(ctx *gin.Context) {
	addrIdStr := ctx.Param("address_id")
	addrId, err := strconv.Atoi(addrIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid address id"})
		return
	}
	var addr models.Address
	if result := dao.DbDao.Model(&models.Address{}).Where("id = ?", addrId).First(&addr); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"address": addr,
	})
}

// AddNewAddress 用户添加新的地址 POST: /api/v1/user/address
func (UserServices) AddNewAddress(ctx *gin.Context) {
	var addNewAddressRequestDto dto.AddNewAddressRequestDto
	if err := ctx.ShouldBindJSON(&addNewAddressRequestDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "提供的数据不合法"})
		return
	}
	if addNewAddressRequestDto.UserId <= 0 || addNewAddressRequestDto.FullName == "" || addNewAddressRequestDto.FullAddress == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "非法请求"})
		return
	}

	newAddress := models.Address{
		UserId:      addNewAddressRequestDto.UserId,
		FullName:    addNewAddressRequestDto.FullName,
		FullAddress: addNewAddressRequestDto.FullAddress,
		PhoneNumber: addNewAddressRequestDto.PhoneNumber,
	}

	if err := dao.DbDao.Create(&newAddress).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "新增地址失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"address": newAddress,
	})
}

// UpdateAddressByUserId 用户修改地址 PUT: /api/v1/user/address/:u_id
func (UserServices) UpdateAddressByUserId(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("u_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}

	var req dto.EditAddressRequestDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "提供的数据不合法"})
		return
	}

	// 查询原始地址，确保属于当前用户
	var addr models.Address
	if err := dao.DbDao.Where("id = ? AND user_id = ?", req.Id, userId).First(&addr).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "地址不存在"})
		return
	}

	// 使用 Updates，只更新需要的字段
	if err := dao.DbDao.Model(&addr).Updates(models.Address{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		FullAddress: req.FullAddress,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "地址更新成功",
		"data":    addr,
	})
}

// DeleteAddressById 用户删除地址 DELETE: /api/v1/user/address/:u_id/:id
func (UserServices) DeleteAddressById(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("u_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}

	addrId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid address id"})
		return
	}

	// 检查是否存在
	var addr models.Address
	if err := dao.DbDao.Where("id = ? AND user_id = ?", addrId, userId).First(&addr).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "地址不存在"})
		return
	}

	// 删除（软删除）
	if err := dao.DbDao.Delete(&addr).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "地址删除成功"})
}
