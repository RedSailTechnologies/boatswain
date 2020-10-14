package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/kraken"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/poseidon"
	krakenRPC "github.com/redsailtechnologies/boatswain/rpc/kraken"
	poseidonRPC "github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "kraken config file path")
	flag.Parse()

	// Kraken
	krakenConfig := &kraken.Config{}
	if err := krakenConfig.YAML(configFile); err != nil {
		logger.Fatal("could not read kraken configuration")
	}

	krakenServer := kraken.New(krakenConfig)
	krakenTwirp := krakenRPC.NewKrakenServer(krakenServer, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))

	// Poseidon
	poseidonConfig := &poseidon.Config{}
	if err := poseidonConfig.YAML(configFile); err != nil {
		logger.Fatal("could not read poseidon configuration")
	}

	poseidonServer := poseidon.New(poseidonConfig)
	poseidonTwirp := poseidonRPC.NewPoseidonServer(poseidonServer, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))

	// Triton
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal("could not read configuration")
	}
	tritonServer := http.FileServer(http.Dir(dir + "/triton"))

	mux := http.NewServeMux()
	mux.Handle(krakenTwirp.PathPrefix(), krakenTwirp)
	mux.Handle(poseidonTwirp.PathPrefix(), poseidonTwirp)
	mux.Handle("/", tritonServer)

	logger.Info("starting leviathan server...ITS HUUUUUUUUUUGE!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":8080", mux))
}
