package main

import (
	"flag"
	"net/http"

	"github.com/redsailtechnologies/boatswain/pkg/logger"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/kraken"
	rpc "github.com/redsailtechnologies/boatswain/rpc/kraken"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "kraken config file path")
	flag.Parse()

	config := &kraken.Config{}
	if err := config.YAML(configFile); err != nil {
		logger.Fatal("could not read configuration")
	}

	server := kraken.New(config)
	twirp := rpc.NewKrakenServer(server, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))
	logger.Info("starting kraken component...RELEASE THE KRAKEN!!!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":8080", twirp))
}
