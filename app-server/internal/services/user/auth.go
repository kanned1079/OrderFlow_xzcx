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
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "用户不存在",
		})
		return
	} else if result.Error != nil {
		services.SendErr500(ctx, result.Error.Error())
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqData.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
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
	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": tokenStr,
	})

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
