package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/redsailtechnologies/boatswain/pkg/kraken"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "kraken config file path")
	flag.Parse()
	config := &kraken.Config{}
	if err := config.YAML(configFile); err != nil {
		log.Fatalf("could not parse configuration")
	}

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	kraken.RegisterKrakenServer(grpcServer, kraken.New(config))
	grpcServer.Serve(listener)
}
