package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendErr500(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"message": message,
	})
}
