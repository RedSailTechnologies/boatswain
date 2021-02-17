package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/agent"
	"github.com/redsailtechnologies/boatswain/pkg/application"
	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/deployment"
	"github.com/redsailtechnologies/boatswain/pkg/git"
	"github.com/redsailtechnologies/boatswain/pkg/health"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	ag "github.com/redsailtechnologies/boatswain/rpc/agent"
	app "github.com/redsailtechnologies/boatswain/rpc/application"
	cl "github.com/redsailtechnologies/boatswain/rpc/cluster"
	dl "github.com/redsailtechnologies/boatswain/rpc/deployment"
	hl "github.com/redsailtechnologies/boatswain/rpc/health"
	rep "github.com/redsailtechnologies/boatswain/rpc/repo"
)

func main() {
	var httpPort, mongoConn, mongoDB string
	flag.StringVar(&httpPort, "http-port", cfg.EnvOrDefaultString("HTTP_PORT", "8080"), "http port")
	flag.StringVar(&mongoConn, "mongo-conn", cfg.EnvOrDefaultString("MONGO_CONNECTION_STRING", ""), "mongodb connection string")
	flag.StringVar(&mongoDB, "mongo-db", cfg.EnvOrDefaultString("MONGO_DB_NAME", "boatswain"), "mongodb database name")
	authCfg := auth.Flags()
	flag.Parse()

	// Storage
	store, err := storage.NewMongo(mongoConn, mongoDB)
	if err != nil {
		logger.Fatal("mongo init failed")
	}

	// Auth
	auth := auth.NewOIDCAgent(authCfg)

	// Twirp Clients
	agentClient := ag.NewAgentActionProtobufClient("http://localhost:8080", &http.Client{}, twirp.WithClientPathPrefix("/agents")) // FIXME

	// Services
	hooks := twirp.ChainHooks(tw.JWTHook(auth), tw.LoggingHooks())

	agent := agent.NewService(store)
	agTwirp := ag.NewAgentServer(agent, tw.LoggingHooks(), twirp.WithServerPathPrefix("/agents"))
	aaTwirp := ag.NewAgentActionServer(agent, tw.LoggingHooks(), twirp.WithServerPathPrefix("/agents"))

	cluster := cluster.NewService(agentClient, auth, store)
	clTwirp := cl.NewClusterServer(cluster, hooks, twirp.WithServerPathPrefix("/api"))

	application := application.NewService(agentClient, auth, store)
	appTwirp := app.NewApplicationServer(application, hooks, twirp.WithServerPathPrefix("/api"))

	repo := repo.NewService(auth, git.DefaultAgent{}, repo.DefaultAgent{}, store)
	repTwirp := rep.NewRepoServer(repo, hooks, twirp.WithServerPathPrefix("/api"))

	deploy := deployment.NewService(agentClient, auth, &git.DefaultAgent{}, store)
	depTwirp := dl.NewDeploymentServer(deploy, hooks, twirp.WithServerPathPrefix("/api"))

	health := health.NewService(application, cluster)
	healthTwirp := hl.NewHealthServer(health, twirp.WithServerPathPrefix("/health"))

	// Browser Client
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal("could not get current directory")
	}
	tritonServer := http.FileServer(http.Dir(dir + "/triton"))

	// Muxing...please stand by...
	mux := http.NewServeMux()
	mux.Handle(agTwirp.PathPrefix(), agTwirp)
	mux.Handle(aaTwirp.PathPrefix(), aaTwirp)
	mux.Handle(appTwirp.PathPrefix(), auth.Wrap(appTwirp))
	mux.Handle(clTwirp.PathPrefix(), auth.Wrap(clTwirp))
	mux.Handle(depTwirp.PathPrefix(), auth.Wrap(depTwirp))
	mux.Handle(repTwirp.PathPrefix(), auth.Wrap(repTwirp))
	mux.Handle(healthTwirp.PathPrefix(), healthTwirp)
	mux.Handle("/", tritonServer) // TODO AdamP - fix multiplexer

	logger.Info("starting leviathan server...ITS HUUUUUUUUUUGE!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, mux))
}
