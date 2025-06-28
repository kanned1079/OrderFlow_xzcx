package routers

import "github.com/gin-gonic/gin"

type GatewayApp struct {
	Id     int32
	Router *gin.Engine
	Mode   string
}

func NewGatewayApp(id int32, ginMode string) *GatewayApp {
	gin.SetMode(ginMode)

	return &GatewayApp{
		Id:     id,
		Router: gin.Default(),
		Mode:   ginMode,
	}
}
