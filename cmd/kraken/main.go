package main

import (
	"flag"

	"github.com/gin-gonic/gin"

	"github.com/redsailtechnologies/boatswain/pkg/kraken"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "kraken config file path")
	flag.Parse()

	router := gin.New()

	// TODO AdamP - improve with better logging/handling
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	r := router.Group("/api/clusters")
	{
		config := new(kraken.ClusterList)
		config.YAML(configFile)
		kraken := kraken.ClustersController{
			Config: config,
		}
		r.GET("", kraken.Clusters)
		r.GET("/:name/namespaces", kraken.Namespaces)
		r.GET("/:name/deployments", kraken.Deployments)
	}

	router.Run(":8081")
}
