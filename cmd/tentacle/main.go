package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	tw "github.com/twitchtv/twirp"

	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/health"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/rpc/agent"
	hl "github.com/redsailtechnologies/boatswain/rpc/health"
)

var (
	client       agent.Agent
	clusterUUID  string
	clusterToken string
	httpPort     string
	timeout      time.Duration
	helmAgent    helm.Agent
	kubeAgent    kube.Agent
)

type healthService struct{}

func (h healthService) Ready() error {
	return nil
}

func init() {
	var bosnURL, timeoutString string
	flag.StringVar(&bosnURL, "bosn-url", cfg.EnvOrDefaultString("BOSN_URL", "http://localhost/"), "boatswain base url")
	flag.StringVar(&clusterUUID, "cluster-uuid", cfg.EnvOrDefaultString("CLUSTER_UUID", ""), "cluster unique id")
	flag.StringVar(&clusterToken, "cluster-token", cfg.EnvOrDefaultString("CLUSTER_TOKEN", ""), "cluster access token")
	flag.StringVar(&httpPort, "http-port", cfg.EnvOrDefaultString("HTTP_PORT", "8080"), "http port")
	flag.StringVar(&timeoutString, "timeout", cfg.EnvOrDefaultString("TIMEOUT", "1s"), "callback timeout to boatswain")
	flag.Parse()

	var err error
	timeout, err = time.ParseDuration(timeoutString)
	if err != nil {
		logger.Fatal("could not parse agent timeout, see https://golang.org/pkg/time/#ParseDuration for details")
	}

	client = agent.NewAgentProtobufClient(bosnURL, &http.Client{}, tw.WithClientPathPrefix("/agents"))

	helmAgent = helm.NewDefaultAgent(restConfig)
	kubeAgent = kube.NewDefaultAgent(k8sConfig)
}

func main() {
	health := health.NewService(&healthService{})
	healthTwirp := hl.NewHealthServer(health, tw.WithServerPathPrefix("/health"))
	mux := http.NewServeMux()
	mux.Handle(healthTwirp.PathPrefix(), healthTwirp)
	go func() {
		logger.Info("starting health checks", "port", httpPort)
		logger.Fatal("health check server failed", "error", http.ListenAndServe(":"+httpPort, mux))
	}()

	logger.Info("starting this kraken tentacle...together they shall rule the world!")
	success := 1
	for {
		time.Sleep(timeout * time.Duration(success))
		actions, err := client.Actions(context.Background(), &agent.ReadActions{
			ClusterUuid:  clusterUUID,
			ClusterToken: clusterToken,
		})
		if err != nil {
			logger.Error("error getting actions", "error", err.Error())
			if time.Duration(success)*timeout < 16*timeout {
				success += success
			}
			continue
		}

		if success > 1 {
			logger.Info("reconnected to boatswain instance")
		}
		success = 1

		if actions.Actions != nil {
			for _, action := range actions.Actions {
				go performAction(action)
			}
		}
	}
}
