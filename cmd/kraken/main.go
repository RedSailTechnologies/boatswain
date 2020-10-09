package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/redsailtechnologies/boatswain/internal/kraken"
	pb "github.com/redsailtechnologies/boatswain/pkg/kraken"
	"github.com/twitchtv/twirp"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "kraken config file path")
	flag.Parse()
	config := &kraken.Config{}
	if err := config.YAML(configFile); err != nil {
		log.Fatalf("could not parse configuration")
	}

	server := kraken.New(config)
	twirp := pb.NewKrakenServer(server, twirp.WithServerPathPrefix("/api"))
	log.Printf(twirp.PathPrefix())
	http.ListenAndServe(":8080", twirp)
}
