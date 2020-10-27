package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/kraken"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	rpc "github.com/redsailtechnologies/boatswain/rpc/kraken"
	"github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "kraken config file path")
	flag.Parse()

	config := &kraken.Config{}
	if err := cfg.YAML(configFile, config); err != nil {
		logger.Fatal("could not read configuration")
	}

	ph := os.Getenv("POSEIDON_SERVICE_HOST")
	pp := os.Getenv("POSEIDON_SERVICE_PORT")
	pe := "http://" + ph + ":" + pp
	poseidon := poseidon.NewPoseidonProtobufClient(pe, &http.Client{}, twirp.WithClientPathPrefix("/api"))

	server := kraken.New(config, poseidon)
	twirp := rpc.NewKrakenServer(server, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))
	logger.Info("starting kraken component...RELEASE THE KRAKEN!!!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":8080", twirp))
}
