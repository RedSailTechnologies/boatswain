package main

import (
	"flag"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/health"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	hl "github.com/redsailtechnologies/boatswain/rpc/health"
	rep "github.com/redsailtechnologies/boatswain/rpc/repo"
)

func main() {
	var httpPort, mongoConn string
	flag.StringVar(&httpPort, "http-port", cfg.EnvOrDefaultString("HTTP_PORT", "8080"), "http port")
	flag.StringVar(&mongoConn, "mongo-conn", cfg.EnvOrDefaultString("MONGO_CONNECTION_STRING", ""), "mongodb connection string")
	authCfg := auth.Flags()
	flag.Parse()

	store, err := storage.NewMongo(mongoConn, "poseidon")
	if err != nil {
		logger.Fatal("mongo init failed")
	}

	authAgent := auth.NewOIDCAgent(authCfg)

	hooks := twirp.ChainHooks(tw.JWTHook(authAgent), tw.LoggingHooks())

	r := repo.NewService(authAgent, git.DefaultAgent{}, helm.DefaultAgent{}, store)
	repTwirp := rep.NewRepoServer(r, hooks, twirp.WithServerPathPrefix("/api"))

	health := health.NewService(r)
	healthTwirp := hl.NewHealthServer(health, twirp.WithServerPathPrefix("/health"))

	mux := http.NewServeMux()
	mux.Handle(repTwirp.PathPrefix(), repTwirp)
	mux.Handle(healthTwirp.PathPrefix(), healthTwirp)

	logger.Info("starting poseidon component...I am Poseidon!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, authAgent.Wrap(mux)))
}
