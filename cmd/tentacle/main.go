package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/redsailtechnologies/boatswain/pkg/cfg"
	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/rpc/agent"
	"github.com/twitchtv/twirp"
	"k8s.io/client-go/kubernetes"
)

var (
	h helm.Agent
	k kube.Agent
)

func main() {
	var bosnURL, clusterUUID, timeoutString string
	flag.StringVar(&bosnURL, "bosn-url", cfg.EnvOrDefaultString("BOSN_URL", "http://localhost/"), "boatswain base url")
	flag.StringVar(&clusterUUID, "cluster-uuid", cfg.EnvOrDefaultString("CLUSTER_UUID", ""), "cluster unique id")
	flag.StringVar(&timeoutString, "timeout", cfg.EnvOrDefaultString("TIMEOUT", "2s"), "callback timeout to boatswain")
	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		logger.Fatal("could not parse agent timeout, see https://golang.org/pkg/time/#ParseDuration for details")
	}

	client := agent.NewAgentProtobufClient(bosnURL, &http.Client{}, twirp.WithClientPathPrefix("/api"))

	registered, err := client.Register(context.TODO(), &agent.RegisterAgent{
		ClusterUuid: clusterUUID,
	})
	if err != nil {
		logger.Fatal("could not register cluster", "error", err)
	}
	clusterToken := registered.ClusterToken

	config, err := kube.GetInClusterConfig()
	if err != nil {
		logger.Fatal("could not get in cluster config, please install in a cluster and give this deployment a service account with sufficient permissions", "error", err)
	}
	k8sConfig, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal("could not convert in cluster config to kubernetes interface", "error", err)
	}

	h = helm.NewDefaultAgent(config)
	k = kube.NewDefaultAgent(k8sConfig)

	for {
		actions, err := client.Actions(context.Background(), &agent.ReadActions{
			ClusterUuid:  clusterUUID,
			ClusterToken: clusterToken,
		})
		if err != nil {
			logger.Error("error getting actions", "error", err)
			continue
		}

		for _, action := range actions.Actions {
			go performAction(client, action)
		}

		time.Sleep(timeout)
	}
}

func performAction(client agent.Agent, action *agent.Action) {
	// TODO - do the thing validate the action and uuid, etc.
	client.Results(context.Background(), &agent.ReturnResult{
		ActionUuid:   action.Uuid,
		ClusterUuid:  action.ClusterUuid,
		ClusterToken: action.ClusterToken,
		Result:       &agent.Result{},
	})
}
