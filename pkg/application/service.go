package application

import (
	"context"

	"github.com/twitchtv/twirp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/application"
	"github.com/redsailtechnologies/boatswain/rpc/cluster"
)

// Service is the implementation of the application service
type Service struct {
	cl  cluster.Cluster
	k8s kube.Agent
}

// NewService returns an initialized instance of the service
func NewService(cluster cluster.Cluster, kube kube.Agent) *Service {
	return &Service{
		cl:  cluster,
		k8s: kube,
	}
}

// All gets all applications currently found in each cluster and their status
func (s Service) All(ctx context.Context, req *pb.ReadApplications) (*pb.ApplicationsRead, error) {
	response := &pb.ApplicationsRead{}

	clusters, err := s.cl.All(context.TODO(), &cluster.ReadClusters{})
	if err != nil {
		logger.Error("error from cluster service", "error", err)
		return nil, err
	}

	for _, c := range clusters.Clusters {
		clientset, err := toClientset(c)
		if err != nil {
			logger.Error("could not get clientset for cluster", "cluster", c.Name)
			return nil, twirp.InternalError("error getting cluster clientset")
		}

		deps, err := s.k8s.GetClusterDeployments(clientset, c.Name)
		if err != nil {
			return nil, twirp.InternalError("error getting deployments for cluster " + c.Name)
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
				c.Name,
				dep.Namespace,
				ready,
			)
		}

		sets, err := s.k8s.GetClusterStatefulSets(clientset, c.Name)
		if err != nil {
			return nil, twirp.InternalError("error getting statefulsets for cluster " + c.Name)
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
				c.Name,
				set.Namespace,
				ready,
			)
		}
	}
	return response, nil
}

func addApplication(resp *pb.ApplicationsRead, name, partOf, version, cluster, namespace string, ready bool) {
	if name == "" || partOf == "" {
		return
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

func toClientset(c *cluster.ClusterRead) (*kubernetes.Clientset, error) {
	restConfig := &rest.Config{
		Host:        c.Endpoint,
		BearerToken: c.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(c.Cert),
		},
	}
	return kubernetes.NewForConfig(restConfig)
}
