package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

// Agent is the interface we use to talk to helm packages
type Agent interface {
	CheckIndex(name, endpoint, token string) bool
	GetChart(name, version, endpoint, token string) ([]byte, error)
	GetChartFromOCIV2(name, version, endpoint, username, password string) ([]byte, error)
}

// DefaultAgent is the default implementation of the Agent interface
type DefaultAgent struct{}

// CheckIndex checks the index.yaml file at the repo's endpoint
func (a DefaultAgent) CheckIndex(name, endpoint, token string) bool {
	r, err := toChartRepo(name, endpoint, token)
	if err != nil {
		return false
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan error, 1)

	// FIXME AdamP - This is really a hack, and leaves goroutines lingering, but the underlying problem
	// is that the DownloadIndexFile call has a timeout, but no context. So while we can see if it took
	// a while to get a response, we can't stop it from taking as long as it wants. My big problem with
	// this approach is that we just leave goroutines out there to finish and we'd rather cancel them.
	// Maybe we should consider just skipping repo status checks altogether, depending on their use.
	go func() {
		_, err := r.DownloadIndexFile()
		if err != nil {
			logger.Warn("error downloading repo index", "error", err)
		}

		select {
		default:
			ch <- err
		case <-ctx.Done():
			return
		}
	}()

	select {
	case err := <-ch:
		return err == nil
	case <-time.After(500 * time.Millisecond):
		return false
	}
}

// GetChart downloads a single chart from a particular chart repo
func (a DefaultAgent) GetChart(name, version, endpoint, token string) ([]byte, error) {
	out := os.TempDir()

	pull := action.NewPull()
	pull.ChartPathOptions = action.ChartPathOptions{
		RepoURL: endpoint,
	}
	pull.Settings = cli.New()
	pull.RepoURL = endpoint
	pull.Version = version

	pull.DestDir = out

	_, err := pull.Run(name)
	if err != nil {
		return nil, err
	}

	path := path.Join(out, fmt.Sprintf("%s-%s.tgz", name, version))
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// TODO - some better dir management here could be good as if this fails we just eventually overload the filesystem
	if err = os.Remove(path); err != nil {
		logger.Warn("could not remove temporary remove temporary file which could lead to disk overusage", "error", err)
	}

	return bytes, nil
}

const (
	v2manifest string = "%s/v2/%s/manifests/%s"
	v2blob     string = "%s/v2/%s/blobs/%s"
)

type manifest struct {
	Layers []struct {
		Digest string `json:"digest"`
	} `json:"layers"`
}

// GetChartFromOCIV2 gets a chart from a V2 compliant OCI registry
func (a DefaultAgent) GetChartFromOCIV2(name, version, endpoint, username, password string) ([]byte, error) {
	u := fmt.Sprintf(v2manifest, endpoint, name, version)
	manifestRequest, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	manifestRequest.SetBasicAuth(username, password)
	manifestRequest.Header.Add("accept", "application/vnd.oci.image.manifest.v1+json")

	client := &http.Client{}
	resp, err := client.Do(manifestRequest)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m := &manifest{}
	json.Unmarshal(b, m)
	if len(m.Layers) != 1 {
		return nil, errors.New("manifest has too many layers")
	}

	u = fmt.Sprintf(v2blob, endpoint, name, m.Layers[0].Digest)
	digestRequest, err := http.NewRequest("GET", u, nil)
	digestRequest.SetBasicAuth(username, password)
	digestRequest.Header.Add("accept", "application/tar+gzip")

	resp, err = client.Do(digestRequest)
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func toChartRepo(name, endpoint, token string) (*repo.ChartRepository, error) {
	providers := []getter.Provider{
		{
			Schemes: []string{"http", "https"},
			New:     getter.NewHTTPGetter,
		},
	}

	// set the username to anything if the token is set
	un := ""
	if token != "" {
		un = "boatswain"
	}

	entry := &repo.Entry{
		Name:     name,
		URL:      endpoint,
		Username: un,
		Password: token,
		// InsecureSkipTLSverify: true, // FIXME - give this option to the user
	}

	return repo.NewChartRepository(entry, providers)
}
