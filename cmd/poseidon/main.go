package main

import (
	"flag"
	"net/http"

	"github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/repo"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	rep "github.com/redsailtechnologies/boatswain/rpc/repo"
)

func main() {
	auth.Flags()
	var httpPort, mongoConn string
	flag.StringVar(&httpPort, "http-port", cfg.EnvOrDefaultString("HTTP_PORT", "8080"), "http port")
	flag.StringVar(&mongoConn, "mongo-conn", cfg.EnvOrDefaultString("MONGO_CONNECTION_STRING", ""), "mongodb connection string")
	flag.Parse()

	store, err := storage.NewMongo(mongoConn, "poseidon")
	if err != nil {
		logger.Fatal("mongo init failed")
	}

	hooks := twirp.ChainHooks(tw.JWTHook(), tw.LoggingHooks())

	repo := repo.NewService(helm.DefaultAgent{}, store)
	repTwirp := rep.NewRepoServer(repo, hooks, twirp.WithServerPathPrefix("/api"))
	logger.Info("starting poseidon component...I am Poseidon!")
	logger.Fatal("server exited", "error", http.ListenAndServe(":"+httpPort, auth.WithJWT(repTwirp)))
}
