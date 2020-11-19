package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/application"
	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	app "github.com/redsailtechnologies/boatswain/rpc/application"
	cl "github.com/redsailtechnologies/boatswain/rpc/cluster"
	rep "github.com/redsailtechnologies/boatswain/rpc/repo"
)

func main() {
	var httpPort, mongoConn string
	flag.StringVar(&httpPort, "http-port", cfg.EnvOrDefaultString("HTTP_PORT", "8080"), "http port")
	flag.StringVar(&mongoConn, "mongo-conn", cfg.EnvOrDefaultString("MONGO_CONNECTION_STRING", ""), "mongodb connection string")
	flag.Parse()

	// Storage
	store, err := storage.NewMongo(mongoConn, "leviathan")
	if err != nil {
		logger.Fatal("mongo init failed")
	}

	// Kraken
	cluster := cluster.NewService(kube.DefaultAgent{}, store)
	clTwirp := cl.NewClusterServer(cluster, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))

	application := application.NewService(cluster, kube.DefaultAgent{})
	appTwirp := app.NewApplicationServer(application, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))

	// Poseidon
	repo := repo.NewService(helm.DefaultAgent{}, store)
	repTwirp := rep.NewRepoServer(repo, logger.TwirpHooks(), twirp.WithServerPathPrefix("/api"))

	// Triton
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal("could not get current directory")
	}
	tritonServer := http.FileServer(http.Dir(dir + "/triton"))

	mux := http.NewServeMux()
	mux.Handle(appTwirp.PathPrefix(), appTwirp)
	mux.Handle(clTwirp.PathPrefix(), clTwirp)
	mux.Handle(repTwirp.PathPrefix(), repTwirp)
	mux.Handle("/", tritonServer) // TODO AdamP - fix multiplexer

	logger.Info("starting leviathan server...ITS HUUUUUUUUUUGE!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, mux))
}
