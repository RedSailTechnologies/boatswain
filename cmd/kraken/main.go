package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	// TODO AdamP - improve with better logging/handling
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	r := router.Group("/api/clusters")
	{
		r.GET("") // list of all clusters
		r.GET("/:name/namespaces")
		r.GET("/:name/deployments") // query param for a ns?
	}

	router.Run()
}
