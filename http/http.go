package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"zpush/conf"
	"zpush/utils"
)

func gatewayHandler(c *gin.Context) {
	servers := utils.GetGatewayServers()
	log.Println(servers)

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data": servers,
	})
}

func StartHTTPServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/gateway", gatewayHandler)
	router.Run(conf.Config().HTTP.Bind)
}
