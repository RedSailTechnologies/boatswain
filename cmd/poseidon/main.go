package main

import (
	"flag"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/poseidon"
	rpc "github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

func main() {
	var configFile, cacheDir string
	flag.StringVar(&configFile, "config", "", "poseidon config file path")
	flag.StringVar(&cacheDir, "cache", "", "poseidon cache path")
	flag.Parse()

	config := &poseidon.Config{}
	if err := cfg.YAML(configFile, config); err != nil {
		logger.Warn("no configuration found or file could not be parsed", "error", err)
	}
	if cacheDir != "" {
		config.CacheDir = cacheDir
	}

	server := poseidon.New(config)
	twirp := rpc.NewPoseidonServer(server, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))
	logger.Info("starting poseidon component...I am Poseidon!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":8080", twirp))
}
