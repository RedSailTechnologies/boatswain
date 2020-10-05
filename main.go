package main

import (
	"github.com/gin-gonic/gin"

	poseidon "poseidon/server"
)

func main() {
	host := gin.Default()
	host.Use(gin.Recovery())

	poseidon.AddRoutes(host)

	host.Run()
}
