package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	appPkg "stay-server/app"
	"stay-server/internal/config"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/utils"
)

func init() {
	config.AppCfg.ReadConfigFile("config/config.yaml")

}

func main() {
	var app = appPkg.NewApp(1, gin.DebugMode)
	fmt.Print(app.InstanceId)
	var count int64
	err := dao.DbDao.Model(&models.User{}).Count(&count).Error
	if err != nil {
		log.Fatal(err)
	}

	var logger utils.Logger

	logger.PrintInfo("count", count)
	logger.PrintWarn("jwt secret: ", config.AppCfg.Runtime.JwtSecret)

	app.GatewayInst.StartApiGateway()
}
