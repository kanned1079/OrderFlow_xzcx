package routers

import (
	"log"
	"stay-server/utils"
)

func (this *GatewayApp) StartApiGateway() {
	var logger utils.Logger
	apiPrefix := this.Router.Group("/api")
	v1 := apiPrefix.Group("/v1")

	this.RegisterPublicRoutes(v1)
	this.RegisterAdminRoutes(v1)
	this.RegisterTraderRoutes(v1)
	this.RegisterUserRoutes(v1)

	logger.PrintSuccess("routers registered successfully.")

	if err := this.Router.Run(":8088"); err != nil {
		log.Println(err)
	}

}
