package http

import (
	"github.com/gin-gonic/gin"
	"zpush/conf"
	"zpush/utils"
)

func gatewayHandler(c *gin.Context) {
	servers := utils.GetGatewayServers()

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    servers,
	})
}

func statsHandler(c *gin.Context) {
	stats := utils.Stats()

	c.JSON(200, gin.H{
		"start_time":    stats.StartTime,
		"total_msg_in":  stats.MsgIn,
		"total_msg_out": stats.MsgOut,
	})
}

func StartHTTPServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/gateway", gatewayHandler)
	router.GET("/stats", statsHandler)

	router.Run(conf.Config().HTTP.Bind)
}
