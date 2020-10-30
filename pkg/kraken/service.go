package kraken

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/twitchtv/twirp"
	"helm.sh/helm/v3/pkg/release"
	"sigs.k8s.io/yaml"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/kraken"
	"github.com/redsailtechnologies/boatswain/rpc/poseidon"
)

// Service is the implementation of the Kraken twirp service
type Service struct {
	clusters  []*Cluster
	kubeAgent kubeAgent
	helmAgent helmAgent
	poseidon  poseidon.Poseidon
}

// New creates the Service with the given configuration
func New(cfg *Config, p poseidon.Poseidon) *Service {
	clusters := make([]*Cluster, len(cfg.Clusters))
	for i, cluster := range cfg.Clusters {
		clusters[i] = &Cluster{
			&pb.Cluster{
				Uuid:     uuid.New().String(),
				Name:     cluster.Name,
				Endpoint: cluster.Endpoint,
				Token:    cluster.Token,
				Cert:     cluster.Cert,
			},
		}
	}
	return &Service{
		clusters:  clusters,
		kubeAgent: defaultKubeAgent{},
		helmAgent: defaultHelmAgent{},
		poseidon:  p,
	}
}

// AddCluster adds a cluster configuration
func (s *Service) AddCluster(ctx context.Context, cluster *pb.Cluster) (*pb.Response, error) {
	s.clusters = append(s.clusters, &Cluster{
		&pb.Cluster{
			Uuid:     uuid.New().String(),
			Name:     cluster.Name,
			Endpoint: cluster.Endpoint,
			Token:    cluster.Token,
			Cert:     cluster.Cert,
		},
	})
	return &pb.Response{}, nil
}

// DeleteCluster deletes a cluster configuration
func (s *Service) DeleteCluster(ctx context.Context, cluster *pb.Cluster) (*pb.Response, error) {
	for i := range s.clusters {
		if s.clusters[i].Uuid == cluster.Uuid {
			s.clusters[i] = s.clusters[len(s.clusters)-1]
			s.clusters = s.clusters[:len(s.clusters)-1]
			return &pb.Response{}, nil
		}
	}
	return nil, twirp.InternalError("cluster not found")
}

// EditCluster edits an existing cluster configuration
func (s *Service) EditCluster(ctx context.Context, cluster *pb.Cluster) (*pb.Response, error) {
	for i := range s.clusters {
		if s.clusters[i].Uuid == cluster.Uuid {
			s.clusters[i] = &Cluster{cluster}
			return &pb.Response{}, nil
		}
	}
	return nil, twirp.InternalError("cluster not found")
}

// Clusters gets all clusters configured
func (s *Service) Clusters(ctx context.Context, req *pb.ClustersRequest) (*pb.ClustersResponse, error) {
	response := &pb.ClustersResponse{
		Clusters: make([]*pb.Cluster, 0),
	}

	for _, cluster := range s.clusters {
		clientset, err := cluster.ToClientset()
		if err != nil {
			logger.Error("could not get clientset for cluster", "cluster", cluster.Name)
			return nil, twirp.InternalError("error getting cluster clientset")
		}

		cluster.Ready = s.kubeAgent.getClusterStatus(clientset, cluster.Name)
		response.Clusters = append(response.Clusters, cluster.Cluster)
	}

	return response, nil
}

// ClusterStatus gets the status of a cluster
func (s *Service) ClusterStatus(ctx context.Context, req *pb.Cluster) (*pb.Cluster, error) {
	cluster := &Cluster{req}
	clientset, err := cluster.ToClientset()
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
		config, err := (&Cluster{cluster}).ToHelmClient("")
		if err != nil {
			logger.Error("could not get helm client for cluster", "cluster", cluster.Name, "error", err)
			return nil, twirp.InternalError("error getting helm client")
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

// UpgradeRelease takes an existing release and updates it
func (s *Service) UpgradeRelease(ctx context.Context, req *pb.UpgradeReleaseRequest) (*pb.Release, error) {
	cluster, err := s.getClusterByName(req.ClusterName)
	if err != nil {
		logger.Error("could not find cluster", "cluster", req.ClusterName, "error", err)
		return nil, twirp.InternalError("cluster not found")
	}

	config, err := cluster.ToHelmClient(req.Namespace)
	if err != nil {
		logger.Error("could not get helm client for cluster", "cluster", req.ClusterName, "error", err)
		return nil, twirp.InternalError("error getting helm client")
	}

	file, err := s.poseidon.DownloadChart(ctx, &poseidon.DownloadRequest{
		ChartName:    req.Chart,
		ChartVersion: req.ChartVersion,
		RepoName:     req.RepoName,
	})
	if err != nil {
		return nil, err
	}

	vals := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(req.Values), &vals); err != nil {
		logger.Error("could not unmarshal values", "values", vals, "error", err)
		return nil, twirp.InternalError("error found in additional values")
	}

	release, err := s.helmAgent.upgradeRelease(config, req.Name, file, req.Namespace, vals)
	if err != nil {
		logger.Error("could not upgrade release", "error", err)
		return nil, twirp.InternalError("error upgrading release")
	}

	return &pb.Release{
		Name:         release.Name,
		Chart:        release.Chart.Metadata.Name,
		Namespace:    release.Namespace,
		AppVersion:   release.Chart.Metadata.AppVersion,
		ChartVersion: release.Chart.Metadata.Version,
		ClusterName:  req.ClusterName,
		Status:       release.Info.Status.String(),
	}, nil
}

func addToReleaseList(list *[]*pb.Releases, search *release.Release, clusterName string) {
	newRelease := &pb.Release{
		Name:         search.Name,
		Chart:        search.Chart.Metadata.Name,
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

func (s *Service) getClusterByName(clusterName string) (*Cluster, error) {
	for _, cluster := range s.clusters {
		if cluster.Name == clusterName {
			return cluster, nil
		}
	}
	return nil, errors.New("cluster not found")
}
