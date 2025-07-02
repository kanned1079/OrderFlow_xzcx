package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/user/dto"
)

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
