package main

import (
	"flag"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/agent"
	"github.com/redsailtechnologies/boatswain/pkg/application"
	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/health"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	ag "github.com/redsailtechnologies/boatswain/rpc/agent"
	app "github.com/redsailtechnologies/boatswain/rpc/application"
	cl "github.com/redsailtechnologies/boatswain/rpc/cluster"
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

	auth := auth.NewOIDCAgent(authCfg)

	hooks := twirp.ChainHooks(tw.JWTHook(auth), tw.LoggingHooks())

	agentClient := ag.NewAgentActionProtobufClient("http://localhost:8080", &http.Client{}, twirp.WithClientPathPrefix("/agents")) // FIXME
	agent := agent.NewService(store)
	agentException := tw.LoggingException{Method: "Actions", Service: "Agent"}
	agTwirp := ag.NewAgentServer(agent, tw.LoggingHooksWithExceptions(agentException), twirp.WithServerPathPrefix("/agents"))
	aaTwirp := ag.NewAgentActionServer(agent, tw.LoggingHooks(), twirp.WithServerPathPrefix("/agents"))

	cluster := cluster.NewService(agentClient, auth, store)
	clTwirp := cl.NewClusterServer(cluster, hooks, twirp.WithServerPathPrefix("/api"))

	application := application.NewService(agentClient, auth, store)
	appTwirp := app.NewApplicationServer(application, hooks, twirp.WithServerPathPrefix("/api"))

	health := health.NewService(application, cluster)
	healthTwirp := hl.NewHealthServer(health, twirp.WithServerPathPrefix("/health"))

	mux := http.NewServeMux()
	mux.Handle(agTwirp.PathPrefix(), agTwirp)
	mux.Handle(aaTwirp.PathPrefix(), aaTwirp)
	mux.Handle(appTwirp.PathPrefix(), auth.Wrap(appTwirp))
	mux.Handle(clTwirp.PathPrefix(), auth.Wrap(clTwirp))
	mux.Handle(healthTwirp.PathPrefix(), healthTwirp)

	logger.Info("starting kraken component...RELEASE THE KRAKEN!!!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, mux))
}
