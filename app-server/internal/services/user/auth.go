package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/internal/services/user/dto"
)

func (UserServices) CheckUserIsValidPreRegister(ctx *gin.Context) {
	var request dto.CheckUserIsValidPreRegisterRequestDto
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "请求格式不合法" + err.Error(),
		})
		return
	}
	// 检查用户美团id
}

func (UserServices) BindWechatAccountAndLogin(ctx *gin.Context) {

}

func (UserServices) Login(ctx *gin.Context) {
	var request dto.UserLoginRequestDto
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "请求格式不合法" + err.Error(),
		})
		return
	}
	log.Println("user login: ", request)
	// 通过临时code请求微信服务器获取用户openid和sessionKey
	// 假设这是从微信获取的
	openid := "ddewfvserhbgrx"
	//sessionKey := "dewfghyjdx"

	var user models.User

	if result := dao.DbDao.Model(&models.User{}).Where(`app_open_id = ?`, openid).First(&user); result.RowsAffected == 0 {
		// 找不到人
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": result.Error.Error(),
		})
		return
	} else if result.Error != nil {
		// 500未知错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.UserLoginResponseDto{
		User: user,
	})
}
