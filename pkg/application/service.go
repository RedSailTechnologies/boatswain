package application

import (
	"context"
	"encoding/json"

	"github.com/redsailtechnologies/boatswain/pkg/auth"
	"github.com/redsailtechnologies/boatswain/pkg/cluster"
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"github.com/redsailtechnologies/boatswain/pkg/storage"
	tw "github.com/redsailtechnologies/boatswain/pkg/twirp"
	"github.com/redsailtechnologies/boatswain/rpc/agent"
	pb "github.com/redsailtechnologies/boatswain/rpc/application"
)

// Service is the implementation of the application service
type Service struct {
	agent agent.AgentAction
	auth  auth.Agent
	cl    *cluster.ReadRepository
	k8s   kube.Agent
}

// NewService returns an initialized instance of the service
func NewService(ag agent.AgentAction, au auth.Agent, s storage.Storage) *Service {
	return &Service{
		agent: ag,
		auth:  au,
		cl:    cluster.NewReadRepository(s),
	}
}

// All gets all applications currently found in each cluster and their status
func (s Service) All(ctx context.Context, req *pb.ReadApplications) (*pb.ApplicationsRead, error) {
	if err := s.auth.Authorize(ctx, auth.Reader); err != nil {
		return nil, tw.ToTwirpError(err, "not authorized")
	}

	response := &pb.ApplicationsRead{}

	clusters, err := s.cl.All()
	if err != nil {
		logger.Error("error from cluster service", "error", err)
		return nil, err
	}

	for _, c := range clusters {
		status, err := s.getClusterStatus(c)
		if err != nil {
			logger.Error("error getting Cluster status", "error", err, "cluster", c.Name())
		}

		if status {
			args := &kube.Args{}
			jsonArgs, err := json.Marshal(args)
			if err != nil {
				logger.Error("error creating agent args", "error", err, "cluster", c.Name())
				continue
			}

			result, err := s.agent.Run(context.Background(), &agent.Action{
				Uuid:           ddd.NewUUID(),
				ClusterUuid:    c.UUID(),
				ClusterToken:   c.Token(),
				ActionType:     agent.ActionType_KUBE_ACTION,
				Action:         string(kube.GetDeployments),
				TimeoutSeconds: 3, // FIXME - configurable?
				Args:           jsonArgs,
			})
			if err != nil {
				logger.Error("error getting Cluster deployments", "error", err, "cluster", c.Name())
				continue
			}

			deps, err := kube.ConvertDeployments(result.Data)
			if err != nil {
				logger.Error("error convertiong luster deployments", "error", "cluster", c.Name())
				continue
			}

			for _, dep := range deps {
				ready := false
				for _, c := range dep.Status.Conditions {
					if c.Type == "Available" && c.Status == "True" {
						ready = true
					}
				}
				addApplication(
					response,
					dep.Labels["app.kubernetes.io/name"],
					dep.Labels["app.kubernetes.io/part-of"],
					dep.Labels["app.kubernetes.io/version"],
					c.Name(),
					dep.Namespace,
					ready,
				)
			}

			result, err = s.agent.Run(context.Background(), &agent.Action{
				Uuid:           ddd.NewUUID(),
				ClusterUuid:    c.UUID(),
				ClusterToken:   c.Token(),
				ActionType:     agent.ActionType_KUBE_ACTION,
				Action:         string(kube.GetStatefulSets),
				TimeoutSeconds: 3, // FIXME - configurable?
				Args:           jsonArgs,
			})
			if err != nil {
				logger.Error("error getting Cluster statefulsets", "error", err, "cluster", c.Name())
				continue
			}

			sets, err := kube.ConvertStatefulSet(result.Data)
			if err != nil {
				logger.Error("error convertiong Cluster statefulsets", "error", "cluster", c.Name())
				continue
			}

			for _, set := range sets {
				ready := false
				for _, c := range set.Status.Conditions {
					if c.Type == "Available" && c.Status == "True" {
						ready = true
					}
				}
				addApplication(
					response,
					set.Labels["app.kubernetes.io/name"],
					set.Labels["app.kubernetes.io/part-of"],
					set.Labels["app.kubernetes.io/version"],
					c.Name(),
					set.Namespace,
					ready,
				)
			}
		} else {
			logger.Warn("could not get applications for cluster as it is not ready", "cluster", c.Name)
		}
	}
	return response, nil
}

// Ready implements the ReadyService method so this service can be part of a health check routine
func (s Service) Ready() error {
	return nil
}

func addApplication(resp *pb.ApplicationsRead, name, partOf, version, cluster, namespace string, ready bool) {
	if name == "" {
		return
	}
	if partOf == "" {
		partOf = "(none)"
	}

	for _, app := range resp.Applications {
		if app.Name == name {
			addApplicationCluster(app, version, cluster, namespace, ready)
			return
		}
	}

	app := &pb.ApplicationRead{
		Name:    name,
		Project: partOf,
	}
	addApplicationCluster(app, version, cluster, namespace, ready)
	resp.Applications = append(resp.Applications, app)
}

func addApplicationCluster(app *pb.ApplicationRead, version, cluster, namespace string, ready bool) {
	if len(app.Clusters) == 0 {
		app.Clusters = make([]*pb.ApplicationCluster, 0)
	}

	for _, c := range app.Clusters {
		if c.ClusterName == cluster {
			return
		}
	}

	cl := &pb.ApplicationCluster{
		ClusterName: cluster,
		Version:     version,
		Namespace:   namespace,
		Ready:       ready,
	}

	app.Clusters = append(app.Clusters, cl)
}

func (s Service) getClusterStatus(c *cluster.Cluster) (bool, error) {
	args := &kube.Args{}
	jsonArgs, err := json.Marshal(args)
	if err != nil {
		return false, err
	}

	result, err := s.agent.Run(context.Background(), &agent.Action{
		Uuid:           ddd.NewUUID(),
		ClusterUuid:    c.UUID(),
		ClusterToken:   c.Token(),
		ActionType:     agent.ActionType_KUBE_ACTION,
		Action:         string(kube.GetStatus),
		TimeoutSeconds: 2, // FIXME - configurable?
		Args:           jsonArgs,
	})
	if err != nil {
		return false, err
	}

	return kube.ConvertStatus(result.Data)
}
