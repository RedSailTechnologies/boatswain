package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/deployment"
	"github.com/redsailtechnologies/boatswain/pkg/health"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	"github.com/redsailtechnologies/boatswain/rpc/cluster"
	dl "github.com/redsailtechnologies/boatswain/rpc/deployment"
	hl "github.com/redsailtechnologies/boatswain/rpc/health"
	"github.com/redsailtechnologies/boatswain/rpc/repo"
)

func main() {
	var httpPort, mongoConn string
	flag.StringVar(&httpPort, "http-port", cfg.EnvOrDefaultString("HTTP_PORT", "8080"), "http port")
	flag.StringVar(&mongoConn, "mongo-conn", cfg.EnvOrDefaultString("MONGO_CONNECTION_STRING", ""), "mongodb connection string")
	authCfg := auth.Flags()
	flag.Parse()

	store, err := storage.NewMongo(mongoConn, "gyarados")
	if err != nil {
		logger.Fatal("mongo init failed")
	}

	authAgent := auth.NewOIDCAgent(authCfg)

	hooks := twirp.ChainHooks(tw.JWTHook(authAgent), tw.LoggingHooks())

	ch := os.Getenv("KRAKEN_SERVICE_HOST")
	cp := os.Getenv("KRAKEN_SERVICE_PORT")
	ce := "http://" + ch + ":" + cp
	clustClient := cluster.NewClusterProtobufClient(ce, &http.Client{}, twirp.WithClientPathPrefix("/api"))

	rh := os.Getenv("POSEIDON_SERVICE_HOST")
	rp := os.Getenv("POSEIDON_SERVICE_PORT")
	re := "http://" + rh + ":" + rp
	repoClient := repo.NewRepoProtobufClient(re, &http.Client{}, twirp.WithClientPathPrefix("/api"))

	deploy := deployment.NewService(authAgent, clustClient, repoClient, store)
	dTwirp := dl.NewDeploymentServer(deploy, hooks, twirp.WithServerPathPrefix("/api"))

	health := health.NewService(deploy)
	healthTwirp := hl.NewHealthServer(health, twirp.WithServerPathPrefix("/health"))

	mux := http.NewServeMux()
	mux.Handle(dTwirp.PathPrefix(), dTwirp)
	mux.Handle(healthTwirp.PathPrefix(), healthTwirp)

	logger.Info("What? MAGIKARP is evolving?")
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, authAgent.Wrap(mux)))
}
