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
	client       agent.Agent
	clusterUUID  string
	clusterToken string
	timeout      time.Duration
	helmAgent    helm.Agent
	kubeAgent    kube.Agent
)

func init() {
	var bosnURL, timeoutString string
	flag.StringVar(&bosnURL, "bosn-url", cfg.EnvOrDefaultString("BOSN_URL", "http://localhost/"), "boatswain base url")
	flag.StringVar(&clusterUUID, "cluster-uuid", cfg.EnvOrDefaultString("CLUSTER_UUID", ""), "cluster unique id")
	flag.StringVar(&timeoutString, "timeout", cfg.EnvOrDefaultString("TIMEOUT", "1s"), "callback timeout to boatswain")

	var err error
	timeout, err = time.ParseDuration(timeoutString)
	if err != nil {
		logger.Fatal("could not parse agent timeout, see https://golang.org/pkg/time/#ParseDuration for details")
	}

	client = agent.NewAgentProtobufClient(bosnURL, &http.Client{}, twirp.WithClientPathPrefix("/agents"))

	config, err := kube.GetInClusterConfig()
	if err != nil {
		logger.Fatal("could not get in cluster config, please install in a cluster and give this deployment a service account with sufficient permissions", "error", err)
	}
	k8sConfig, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal("could not convert in cluster config to kubernetes interface", "error", err)
	}

	helmAgent = helm.NewDefaultAgent(config)
	kubeAgent = kube.NewDefaultAgent(k8sConfig)
}

func main() {
	// TODO AdamP - health checks?
	registered := false
	for !registered {
		result, err := client.Register(context.TODO(), &agent.RegisterAgent{
			ClusterUuid: clusterUUID,
		})
		if err != nil {
			logger.Warn("could not register cluster", "error", err)
			time.Sleep(5 * time.Second) // TODO AdamP - lets scale this up as failures increase
		} else {
			registered = true
			clusterToken = result.ClusterToken
			logger.Info("cluster registered")
		}
	}

	for {
		time.Sleep(timeout) // TODO AdamP - lets scale this up as failures increase
		actions, err := client.Actions(context.Background(), &agent.ReadActions{
			ClusterUuid:  clusterUUID,
			ClusterToken: clusterToken,
		})
		if err != nil {
			logger.Error("error getting actions", "error", err)
			continue // TODO AdamP - lets scale this up as failures increase
		}

		for _, action := range actions.Actions {
			go performAction(action)
		}
	}
}
