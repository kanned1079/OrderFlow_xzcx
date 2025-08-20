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

func (this *TraderServices) FetchAddressInfoByAddressId(ctx *gin.Context) {
	addressIdStr := ctx.Param("a_id")
	addressId, err := strconv.ParseInt(addressIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "地址ID格式错误"})
		return
	}

	var existingAddress models.Address
	if result := dao.DbDao.Model(&models.Address{}).Where("id = ?", addressId).Unscoped().First(&existingAddress); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "找不到指定Id的地址信息",
		})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询出错 " + result.Error.Error(),
		})
		return
	}

	var response dto.GetOrderInfoByAddressIdResponseDto = dto.GetOrderInfoByAddressIdResponseDto{
		Address: existingAddress,
		Message: "success",
	}

	ctx.JSON(http.StatusOK, response)
}
