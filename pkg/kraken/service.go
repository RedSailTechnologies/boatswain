package kraken

import (
	"context"
	"errors"

	"github.com/twitchtv/twirp"
	"helm.sh/helm/v3/pkg/release"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/kraken"
)

// Service is the implementation of the Kraken twirp service
type Service struct {
	config    *Config
	kubeAgent kubeAgent
	helmAgent helmAgent
}

// New creates the Service with the given configuration
func New(c *Config) *Service {
	return &Service{
		config:    c,
		kubeAgent: defaultKubeAgent{},
		helmAgent: defaultHelmAgent{},
	}
}

// Clusters gets all clusters configured
func (s *Service) Clusters(ctx context.Context, req *pb.ClustersRequest) (*pb.ClustersResponse, error) {
	response := &pb.ClustersResponse{
		Clusters: make([]*pb.Cluster, 0),
	}

	for _, cluster := range s.config.Clusters {
		clientset, err := s.config.ToClientset(cluster.Name)
		if err != nil {
			logger.Error("could not get clientset for cluster", "cluster", cluster.Name)
			return nil, twirp.InternalError("error getting cluster clientset")
		}

		response.Clusters = append(response.Clusters, &pb.Cluster{
			Name:     cluster.Name,
			Endpoint: cluster.Endpoint,
			Ready:    s.kubeAgent.getClusterStatus(clientset, cluster.Name),
		})
	}

	return response, nil
}

// ClusterStatus gets the status of a cluster
func (s *Service) ClusterStatus(ctx context.Context, cluster *pb.Cluster) (*pb.Cluster, error) {
	clientset, err := s.config.ToClientset(cluster.Name)
	if err != nil {
		logger.Error("could not get clientset for cluster", "cluster", cluster.Name)
		return nil, twirp.InternalError("error getting cluster clientset")
	}

	return &pb.Cluster{
		Name:     cluster.Name,
		Endpoint: cluster.Endpoint,
		Ready:    s.kubeAgent.getClusterStatus(clientset, cluster.Name),
	}, nil
}

// Releases gets all releases based on the clusters in the request
func (s *Service) Releases(ctx context.Context, req *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {
	releaseList := make([]*pb.Releases, 0)
	for _, cluster := range req.Clusters {
		config, err := s.config.ToHelmClient(cluster.Name)
		if err != nil {
			logger.Error("could not get clientset for cluster", "cluster", cluster.Name, "error", err)
			return nil, twirp.InternalError("error getting cluster clientset")
		}

		// get releases for cluster
		releases, err := s.helmAgent.getReleases(config, cluster.Name)
		if err != nil {
			logger.Error("could not get releases for cluster", "cluster", cluster.Name, "error", err)
			continue
		}

		for _, release := range releases {
			addToReleaseList(&releaseList, release, cluster.Name)
		}
	}
	return &pb.ReleaseResponse{ReleaseLists: releaseList}, nil
}

// ReleaseStatus gets an updated status for a particular release
func (s *Service) ReleaseStatus(ctx context.Context, release *pb.Release) (*pb.Release, error) {
	return nil, errors.New("not implemented")
}

func addToReleaseList(list *[]*pb.Releases, search *release.Release, clusterName string) {
	newRelease := &pb.Release{
		Namespace:    search.Namespace,
		AppVersion:   search.Chart.Metadata.AppVersion,
		ChartVersion: search.Chart.Metadata.Version,
		ClusterName:  clusterName,
		Status:       search.Info.Status.String(),
	}

	for _, release := range *list {
		if release.Name == search.Name && release.Chart == search.Chart.Metadata.Name {
			release.Releases = append(release.Releases, newRelease)
			return
		}
	}

	*list = append(*list, &pb.Releases{
		Name:  search.Name,
		Chart: search.Chart.Metadata.Name,
		Releases: []*pb.Release{
			newRelease,
		},
	})
}
