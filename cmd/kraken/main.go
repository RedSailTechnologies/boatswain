package main

import (
	"github.com/gin-gonic/gin"

	"github.com/redsailtechnologies/boatswain/pkg/kraken"
)

func main() {
	router := gin.New()

	// TODO AdamP - improve with better logging/handling
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	r := router.Group("/api/clusters")
	{
		kraken := new(kraken.ClustersController)
		r.GET("", kraken.Clusters)
		r.GET("/:name/namespaces", kraken.Namespaces)
		r.GET("/:name/deployments", kraken.Deployments)
	}

	router.Run()
}
