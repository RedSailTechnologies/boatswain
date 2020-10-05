package main

import (
	"github.com/gin-gonic/gin"

	"github.com/redsailtechnologies/boatswain/poseidon/server"
)

func main() {
	host := gin.Default()
	host.Use(gin.Recovery())

	server.AddRoutes(host)

	host.Run()
}
