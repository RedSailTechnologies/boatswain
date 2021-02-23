package main

import (
	"flag"
	"net/http"
	"os"

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
	tr "github.com/redsailtechnologies/boatswain/rpc/trigger"
)

func main() {
	var actionEndpoint, httpPort, mongoConn, mongoDB string
	flag.StringVar(&actionEndpoint, "action-endpoint", cfg.EnvOrDefaultString("ACTION_ENDPOINT", "http://localhost:8080"), "agent action service endpoint")
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
	actionClient := ag.NewAgentActionProtobufClient("http://localhost:8080", &http.Client{}, twirp.WithClientPathPrefix("/agents"))

	// Services
	hooks := twirp.ChainHooks(tw.JWTHook(auth), tw.LoggingHooks())
	opts := twirp.WithServerPathPrefix("/api")

	agent := agent.NewService(store)
	agentException := tw.LoggingException{Method: "Actions", Service: "Agent"}
	agTwirp := ag.NewAgentServer(agent, tw.LoggingHooksWithExceptions(agentException), twirp.WithServerPathPrefix("/agents"))
	aaTwirp := ag.NewAgentActionServer(agent, tw.LoggingHooks(), twirp.WithServerPathPrefix("/agents"))

	cluster := cluster.NewService(actionClient, auth, store)
	clTwirp := cl.NewClusterServer(cluster, hooks, opts)

	application := application.NewService(actionClient, auth, store)
	appTwirp := app.NewApplicationServer(application, hooks, opts)

	repo := repo.NewService(auth, git.DefaultAgent{}, repo.DefaultAgent{}, store)
	repTwirp := rep.NewRepoServer(repo, hooks, opts)

	deploy := deployment.NewService(actionClient, auth, &git.DefaultAgent{}, store)
	depTwirp := dl.NewDeploymentServer(deploy, hooks, opts)
	trigTwirp := tr.NewTriggerServer(deploy, tw.LoggingHooks(), opts)

	health := health.NewService(agent, application, cluster, deploy, repo)
	healthTwirp := hl.NewHealthServer(health, twirp.WithServerPathPrefix("/health"))

	// Muxing...please stand by...
	mux := http.NewServeMux()
	mux.Handle(agTwirp.PathPrefix(), agTwirp)
	mux.Handle(aaTwirp.PathPrefix(), aaTwirp)
	mux.Handle(appTwirp.PathPrefix(), auth.Wrap(appTwirp))
	mux.Handle(clTwirp.PathPrefix(), auth.Wrap(clTwirp))
	mux.Handle(depTwirp.PathPrefix(), auth.Wrap(depTwirp))
	mux.Handle(repTwirp.PathPrefix(), auth.Wrap(repTwirp))
	mux.Handle(trigTwirp.PathPrefix(), auth.Wrap(trigTwirp))
	mux.Handle(healthTwirp.PathPrefix(), healthTwirp)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat("triton/" + r.RequestURI); os.IsNotExist(err) {
			http.ServeFile(w, r, "triton/index.html")
		}
		http.ServeFile(w, r, "triton/"+r.RequestURI)
	}) // TODO AdamP - fix multiplexer

	logger.Info("starting leviathan server...ITS HUUUUUUUUUUGE!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, mux))
}
