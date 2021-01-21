package main

import (
	"flag"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/deployment"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/health"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	dl "github.com/redsailtechnologies/boatswain/rpc/deployment"
	hl "github.com/redsailtechnologies/boatswain/rpc/health"
)

func main() {
	var httpPort, mongoConn, mongoDB string
	flag.StringVar(&httpPort, "http-port", cfg.EnvOrDefaultString("HTTP_PORT", "8080"), "http port")
	flag.StringVar(&mongoConn, "mongo-conn", cfg.EnvOrDefaultString("MONGO_CONNECTION_STRING", ""), "mongodb connection string")
	flag.StringVar(&mongoDB, "mongo-db", cfg.EnvOrDefaultString("MONGO_DB_NAME", "boatswain"), "mongodb database name")
	authCfg := auth.Flags()
	flag.Parse()

	store, err := storage.NewMongo(mongoConn, mongoDB)
	if err != nil {
		logger.Fatal("mongo init failed")
	}

	authAgent := auth.NewOIDCAgent(authCfg)

	hooks := twirp.ChainHooks(tw.JWTHook(authAgent), tw.LoggingHooks())

	deploy := deployment.NewService(authAgent, &git.DefaultAgent{}, store)
	dTwirp := dl.NewDeploymentServer(deploy, hooks, twirp.WithServerPathPrefix("/api"))

	health := health.NewService(deploy)
	healthTwirp := hl.NewHealthServer(health, twirp.WithServerPathPrefix("/health"))

	mux := http.NewServeMux()
	mux.Handle(dTwirp.PathPrefix(), dTwirp)
	mux.Handle(healthTwirp.PathPrefix(), healthTwirp)

	logger.Info("What? MAGIKARP is evolving?")
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, authAgent.Wrap(mux)))
}
