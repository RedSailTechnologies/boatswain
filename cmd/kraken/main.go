package main

import (
	"flag"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/kraken"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	rpc "github.com/redsailtechnologies/boatswain/rpc/kraken"
	"github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

func main() {
	var poseidonHost, poseidonPort string
	flag.StringVar(&poseidonHost, "poseidon-host", cfg.EnvOrDefaultString("POSEIDON_SERVICE_HOST", "not.found"), "poseidon service host")
	flag.StringVar(&poseidonPort, "poseidon-port", cfg.EnvOrDefaultString("POSEIDON_SERVICE_PORT", "0000"), "poseidon service host")
	flag.Parse()

	poseidonEndpoint := "http://" + poseidonHost + ":" + poseidonPort
	poseidon := poseidon.NewPoseidonProtobufClient(poseidonEndpoint, &http.Client{}, twirp.WithClientPathPrefix("/api"))

	server := kraken.New(config, poseidon)
	twirp := rpc.NewKrakenServer(server, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))
	logger.Info("starting kraken component...RELEASE THE KRAKEN!!!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":8080", twirp))
}
