package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
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

	grpcListener, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	httpListener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	kraken.RegisterKrakenServer(grpcServer, kraken.New(config))

	wrappedGrpc := grpcweb.WrapServer(grpcServer)
	httpServer := http.Server{}
	httpServer.Handler = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
		} else {
			resp.WriteHeader(404)
		}
	})

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		log.Printf("starting grpc server")
		grpcServer.Serve(grpcListener)
		log.Printf("grpc server exited")
		wg.Done()
	}()
	go func() {
		log.Printf("starting http server")
		httpServer.Serve(httpListener)
		log.Printf("http server exited")
		wg.Done()
	}()
	wg.Wait()
}
