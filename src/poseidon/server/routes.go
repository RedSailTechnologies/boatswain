package server

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(server *gin.Engine) {
	server.GET("/api/clusters", func(c *gin.Context) { c.Data(200, gin.MIMEPlain, nil) })
}
