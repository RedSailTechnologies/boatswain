package main

import (
	"flag"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/application"
	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	app "github.com/redsailtechnologies/boatswain/rpc/application"
	cl "github.com/redsailtechnologies/boatswain/rpc/cluster"
)

func main() {
	var httpPort, mongoConn string
	flag.StringVar(&httpPort, "http-port", cfg.EnvOrDefaultString("HTTP_PORT", "8080"), "http port")
	flag.StringVar(&mongoConn, "mongo-conn", cfg.EnvOrDefaultString("MONGO_CONNECTION_STRING", ""), "mongodb connection string")
	flag.Parse()

	store, err := storage.NewMongo(mongoConn, "kraken")
	if err != nil {
		logger.Fatal("mongo init failed")
	}

	cluster := cluster.NewService(kube.DefaultAgent{}, store)
	clTwirp := cl.NewClusterServer(cluster, tw.LoggingHooks(), twirp.WithServerPathPrefix("/api"))

	application := application.NewService(cluster, kube.DefaultAgent{})
	appTwirp := app.NewApplicationServer(application, tw.LoggingHooks(), twirp.WithServerPathPrefix("/api"))

	mux := http.NewServeMux()
	mux.Handle(appTwirp.PathPrefix(), appTwirp)
	mux.Handle(clTwirp.PathPrefix(), clTwirp)

	logger.Info("starting kraken component...RELEASE THE KRAKEN!!!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, mux))
}
