package main

import (
	"flag"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/application"
	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/health"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	app "github.com/redsailtechnologies/boatswain/rpc/application"
	cl "github.com/redsailtechnologies/boatswain/rpc/cluster"
	hl "github.com/redsailtechnologies/boatswain/rpc/health"
)

func main() {
	var httpPort, mongoConn string
	flag.StringVar(&httpPort, "http-port", cfg.EnvOrDefaultString("HTTP_PORT", "8080"), "http port")
	flag.StringVar(&mongoConn, "mongo-conn", cfg.EnvOrDefaultString("MONGO_CONNECTION_STRING", ""), "mongodb connection string")
	authCfg := auth.Flags()
	flag.Parse()

	store, err := storage.NewMongo(mongoConn, "kraken")
	if err != nil {
		logger.Fatal("mongo init failed")
	}

	authAgent := auth.NewOIDCAgent(authCfg)

	hooks := twirp.ChainHooks(tw.JWTHook(authAgent), tw.LoggingHooks())

	cluster := cluster.NewService(authAgent, kube.DefaultAgent{}, store)
	clTwirp := cl.NewClusterServer(cluster, hooks, twirp.WithServerPathPrefix("/api"))

	application := application.NewService(authAgent, cluster, kube.DefaultAgent{})
	appTwirp := app.NewApplicationServer(application, hooks, twirp.WithServerPathPrefix("/api"))

	health := health.NewService(application, cluster)
	healthTwirp := hl.NewHealthServer(health, twirp.WithServerPathPrefix("/health"))

	mux := http.NewServeMux()
	mux.Handle(appTwirp.PathPrefix(), appTwirp)
	mux.Handle(clTwirp.PathPrefix(), clTwirp)
	mux.Handle(healthTwirp.PathPrefix(), healthTwirp)

	logger.Info("starting kraken component...RELEASE THE KRAKEN!!!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, authAgent.Wrap(mux)))
}
