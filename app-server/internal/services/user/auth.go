package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services"
	"stay-server/internal/services/user/dto"
	"strconv"
	"time"
)

func (this *UserServices) Login(ctx *gin.Context) {
	//u := utils.Utils{}
	var reqData dto.UserLoginRequestDto
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "请求格式不合法" + err.Error(),
		})
		return
	}
	//time.Sleep(time.Microsecond * 100)
	var user models.User
	if result := dao.DbDao.Model(&models.User{}).Where("phone_number = ?", reqData.PhoneNumber).First(&user); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 用户没找到
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "用户不存在",
		})
		return
	} else if result.Error != nil {
		services.SendErr500(ctx, result.Error.Error())
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqData.Password))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "密码错误",
		})
		return
	}

	tokenStr, err := this.utils.GenerateAccessToken(user)
	if err != nil {
		services.SendErr500(ctx, err.Error())
		return
	}

	now := time.Now()
	if err := dao.DbDao.Model(&user).Update("last_login_at", now).Error; err != nil {
		services.SendErr500(ctx, "更新登录时间失败: "+err.Error())
		return
	}

	user.Password = ""

	if user.Role == "trader" {
		var merchant models.Merchant
		if result := dao.DbDao.Model(&models.Merchant{}).Where("user_id = ?", user.Id).First(&merchant); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "加载商家信息时出现错误 " + result.Error.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"user":     user,
			"merchant": merchant,
			"token":    tokenStr,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"user":  user,
			"token": tokenStr,
		})
	}

}

func (this *UserServices) Register(ctx *gin.Context) {
	//u := utils.Utils{}
	var reqData dto.UserRegisterRequestDto
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "请求格式不合法" + err.Error(),
		})
		return
	}

	var count int64
	dao.DbDao.Model(&models.User{}).Where("phone_number = ?", reqData.PhoneNumber).Count(&count)
	if count > 0 {
		services.SendErr500(ctx, "该手机号已注册")
		return
	}
	var newUser models.User = models.User{
		Username:    reqData.Username,
		PhoneNumber: reqData.PhoneNumber,
		Status:      true,
		Role:        "user",
	}

	// todo 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), bcrypt.DefaultCost)
	if err != nil {
		services.SendErr500(ctx, "密码加密失败: "+err.Error())
		return
	}
	newUser.Password = string(hashedPassword)

	if result := dao.DbDao.Create(&newUser); result.Error != nil {
		services.SendErr500(ctx, result.Error.Error())
		return
	}

	tokenStr, err := this.utils.GenerateAccessToken(newUser)
	if err != nil {
		services.SendErr500(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  newUser,
		"token": tokenStr,
	})
}

// UpdateUserPassword 用户修改密码 PATCH: /api/v1/user/:u_id/password/update
func (this *UserServices) UpdateUserPassword(ctx *gin.Context) {
	userIdStr := ctx.Param("u_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "非法用户Id: " + err.Error(),
		})
		return
	}

	this.utils.Logger.PrintWarn(userId)
	var reqData dto.UserUpdatePasswordRequestDto
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "请求格式不合法: " + err.Error(),
		})
		return
	}

	this.utils.Logger.PrintWarn(reqData)
	// 查找用户
	var user models.User
	if result := dao.DbDao.First(&user, userId); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "用户不存在",
			})
			return
		}
		services.SendErr500(ctx, result.Error.Error())
		return
	}

	// 校验旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqData.PreviousPassword)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "旧密码错误",
		})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		services.SendErr500(ctx, "密码加密失败: "+err.Error())
		return
	}

	// 更新数据库
	if err := dao.DbDao.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		services.SendErr500(ctx, "更新密码失败: "+err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "密码更新成功",
	})
}
