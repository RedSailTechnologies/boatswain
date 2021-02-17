package main

import (
	"context"
	"encoding/json"

	"github.com/redsailtechnologies/boatswain/pkg/helm"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/rpc/agent"
)

func performAction(action *agent.Action) {
	logger.Info("performing action", "id", action.Uuid, "action", action.Action, "type", action.ActionType)
	response := &agent.ReturnResult{
		ActionUuid:   action.Uuid,
		ClusterUuid:  action.ClusterUuid,
		ClusterToken: action.ClusterToken,
		Result:       &agent.Result{},
	}

	if action.ClusterUuid != clusterUUID || action.ClusterToken != clusterToken {
		response.Result.Error = "invalid cluster uuid or token"
		client.Results(context.Background(), nil)
		return
	}

	switch action.ActionType {
	case agent.ActionType_KUBE_ACTION:
		client.Results(context.Background(), runKubeAction(action, response))
	case agent.ActionType_HELM_ACTION:
		client.Results(context.Background(), runHelmAction(action, response))
	default:
		response.Result.Error = "action type not found"
		client.Results(context.Background(), response)
	}
	logger.Info("result sent", "id", action.Uuid, "error", response.Result.Error)
}

func runKubeAction(action *agent.Action, response *agent.ReturnResult) *agent.ReturnResult {
	args, err := kube.ConvertArgs(action.Args)
	if err != nil {
		response.Result.Error = err.Error()
		return response
	}

	var result *kube.Result
	switch kube.AgentAction(action.Action) {
	case kube.GetStatus:
		result, err = kubeAgent.GetStatus(args)
	case kube.GetDeployments:
		result, err = kubeAgent.GetDeployments(args)
	case kube.GetStatefulSets:
		result, err = kubeAgent.GetStatefulSets(args)
	default:
		response.Result.Error = "kube agent action not found"
		return response
	}

	if err != nil {
		response.Result.Error = err.Error()
		return response
	}

	data, err := json.Marshal(result)
	if err != nil {
		response.Result.Error = err.Error()
		return response
	}
	response.Result.Data = data
	return response
}

func runHelmAction(action *agent.Action, response *agent.ReturnResult) *agent.ReturnResult {
	args, err := helm.ConvertArgs(action.Args)
	if err != nil {
		response.Result.Error = err.Error()
		return response
	}

	var result *helm.Result
	switch helm.AgentAction(action.Action) {
	case helm.Install:
		result, err = helmAgent.Install(args)
	case helm.Rollback:
		result, err = helmAgent.Rollback(args)
	case helm.Test:
		result, err = helmAgent.Test(args)
	case helm.Uninstall:
		result, err = helmAgent.Uninstall(args)
	case helm.Upgrade:
		result, err = helmAgent.Upgrade(args)
	default:
		response.Result.Error = "helm agent action not found"
		return response
	}

	if err != nil {
		response.Result.Error = err.Error()
		return response
	}

	data, err := json.Marshal(result)
	if err != nil {
		response.Result.Error = err.Error()
		return response
	}
	response.Result.Data = data
	return response
}
