package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"stay-server/internal/middlewares"
	"stay-server/utils"
	"time"
)

func (this *GatewayApp) StartApiGateway() {
	var logger utils.Logger
	apiPrefix := this.Router.Group("/api")
	v1 := apiPrefix.Group("/v1")

	this.RegisterPublicRoutes(v1)
	this.RegisterAdminRoutes(v1)
	this.RegisterTraderRoutes(v1)
	this.RegisterUserRoutes(v1)

	// 注册静态文件访问路径
	this.Router.Static("/static", "./uploads")

	v1.POST("/file/upload", middlewares.RequireAuth(), func(ctx *gin.Context) {
		// 从 multipart/form-data 中获取文件字段 "file"
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "上传文件失败: " + err.Error()})
			return
		}

		// 创建 uploads 目录（如果不存在）
		savePath := "./uploads"
		_ = os.MkdirAll(savePath, os.ModePerm)

		// 为避免文件名冲突，可以用时间戳或 UUID 命名
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		fullPath := filepath.Join(savePath, filename)

		// 保存文件到服务器本地
		if err := ctx.SaveUploadedFile(file, fullPath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "保存文件失败: " + err.Error()})
			return
		}

		// 构造可公开访问的 URL
		publicURL := "/static/" + filename

		ctx.JSON(http.StatusOK, gin.H{
			"message":  "上传成功",
			"filename": filename,
			"url":      publicURL,
		})
	})

	logger.PrintSuccess("routers registered successfully.")

	if err := this.Router.Run(":8088"); err != nil {
		log.Println(err)
	}

}
