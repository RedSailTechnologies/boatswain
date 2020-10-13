package main

import (
	"flag"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/poseidon"
	rpc "github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "poseidon config file path")
	flag.Parse()

	config := &poseidon.Config{}
	if err := config.YAML(configFile); err != nil {
		logger.Fatal("could not read configuration")
	}

	server := poseidon.New(config)
	twirp := rpc.NewPoseidonServer(server, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))
	logger.Info("starting poseidon component...I am Poseidon!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":8080", twirp))
}
