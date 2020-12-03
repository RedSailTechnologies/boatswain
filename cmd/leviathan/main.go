package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/application"
	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	app "github.com/redsailtechnologies/boatswain/rpc/application"
	cl "github.com/redsailtechnologies/boatswain/rpc/cluster"
	rep "github.com/redsailtechnologies/boatswain/rpc/repo"
)

func main() {
	var httpPort, mongoConn string
	flag.StringVar(&httpPort, "http-port", cfg.EnvOrDefaultString("HTTP_PORT", "8080"), "http port")
	flag.StringVar(&mongoConn, "mongo-conn", cfg.EnvOrDefaultString("MONGO_CONNECTION_STRING", ""), "mongodb connection string")
	authCfg := auth.Flags()
	flag.Parse()

	// Storage
	store, err := storage.NewMongo(mongoConn, "leviathan")
	if err != nil {
		logger.Fatal("mongo init failed")
	}

	// Auth
	authAgent := auth.NewOIDCAgent(authCfg)

	// Services
	hooks := twirp.ChainHooks(tw.JWTHook(authAgent), tw.LoggingHooks())

	cluster := cluster.NewService(authAgent, kube.DefaultAgent{}, store)
	clTwirp := cl.NewClusterServer(cluster, hooks, twirp.WithServerPathPrefix("/api"))

	application := application.NewService(authAgent, cluster, kube.DefaultAgent{})
	appTwirp := app.NewApplicationServer(application, hooks, twirp.WithServerPathPrefix("/api"))

	repo := repo.NewService(authAgent, helm.DefaultAgent{}, store)
	repTwirp := rep.NewRepoServer(repo, hooks, twirp.WithServerPathPrefix("/api"))

	// Client
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
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, authAgent.Wrap(mux)))
}
