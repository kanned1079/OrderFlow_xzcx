package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"runtime"
	appPkg "stay-server/app"
	"stay-server/internal/config"
	"stay-server/internal/dao"
	"stay-server/internal/models"
	"stay-server/utils"
	"strings"
	"sync"
	"time"
)

func init() {
	showSystemConfig()
	config.AppCfg.ReadConfigFile("config/config.yaml")
}

func checkAdminExist() {
	var myLogger utils.Logger
	var count int64
	err := dao.DbDao.Model(&models.User{}).Count(&count).Error
	if err != nil {
		log.Panicln(err)
	}
	if count <= 0 {
		myLogger.PrintWarn("你需要先设置一个管理员")
		var phoneNumber, inputPassword string

		fmt.Print("请输入管理员手机号: ")
		if _, err := fmt.Scanln(&phoneNumber); err != nil {
			log.Panicln("手机号读取失败:", err)
		}
		phoneNumber = strings.TrimSpace(phoneNumber)
		if len(phoneNumber) != 11 {
			log.Panicln("手机号必须为 11 位数字")
		}

		fmt.Print("请输入管理员密码: ")
		if _, err := fmt.Scanln(&inputPassword); err != nil {
			log.Panicln("密码读取失败:", err)
		}
		inputPassword = strings.TrimSpace(inputPassword)
		if len(inputPassword) < 6 {
			log.Panicln("密码长度不能小于 6 位")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Panicln("密码加密失败:", err)
		}

		if err := dao.DbDao.Create(&models.User{
			Username:    "sysadmin",
			Role:        "admin",
			Status:      true,
			PhoneNumber: phoneNumber,
			Password:    string(hashedPassword),
		}).Error; err != nil {
			log.Panicln("创建管理员失败:", err)
		}

		myLogger.PrintInfo("你设置的是", phoneNumber)

	} else {
		myLogger.PrintSuccess("已设置管理员 Starting gateway api services...")
	}
}

func main() {
	var app = appPkg.NewApp(1, gin.DebugMode)
	checkAdminExist()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.GatewayInst.StartApiGateway()
	}()

	wg.Wait()
}

func showSystemConfig() {
	var myLogger utils.Logger
	myLogger.PrintInfo(fmt.Sprintf("OS(Arch): %s %s", runtime.GOOS, runtime.GOARCH))
	myLogger.PrintInfo(fmt.Sprintf("CPU(s): %v", runtime.NumCPU()))
	myLogger.PrintInfo(fmt.Sprintf("当前时间: %v", time.Now().Local()))
}
