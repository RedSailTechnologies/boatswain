package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

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

	krakenServer := kraken.New(config)
	krakenTwirp := rpc.NewKrakenServer(krakenServer, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal("could not read configuration")
	}
	tritonServer := http.FileServer(http.Dir(dir + "/triton"))

	mux := http.NewServeMux()
	mux.Handle(krakenTwirp.PathPrefix(), krakenTwirp)
	mux.Handle("/", tritonServer)

	logger.Info("starting leviathan server...ITS HUUUUUUUUUUGE!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":8080", mux))
}
