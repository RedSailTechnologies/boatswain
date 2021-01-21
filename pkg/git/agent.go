package git

import (
	"io/ioutil"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
)

// Agent is the interface for git interaction
type Agent interface {
	CheckRepo(endpoint, username, password string) bool
	GetFile(endpoint, branch, path, username, password string) []byte
}

// DefaultAgent is the default implementaiton of the Agent interface
type DefaultAgent struct{}

// CheckRepo tries to fetch from the endpoint to check if it is valid
func (a DefaultAgent) CheckRepo(endpoint, username, password string) bool {
	store := memory.NewStorage()
	rc := config.RemoteConfig{
		URLs: []string{
			endpoint,
		},
	}
	r := git.NewRemote(store, &rc)

	fo := git.FetchOptions{
		Auth: getAuth(username, password),
		RefSpecs: []config.RefSpec{
			"+refs/heads/*:refs/remotes/origin/*",
		},
	}
	err := r.Fetch(&fo)

	// check that in addition to not having an error, we have at least one ref
	if err != nil || len(store.ReferenceStorage) == 0 {
		return false
	}
	return true
}

// GetFile gets a single file from the repo
func (a DefaultAgent) GetFile(endpoint, branch, path, username, password string) []byte {
	store := memory.NewStorage()
	fs := memfs.New()

	_, err := git.Clone(store, fs, &git.CloneOptions{
		URL:           endpoint,
		Auth:          getAuth(username, password),
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		logger.Error("error cloning git repo", "error", err)
		return nil
	}

	f, err := fs.Open(path)
	if err != nil {
		logger.Error("could not find specified file", "error", err)
		return nil
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Error("error reading file", "error", err)
		return nil
	}
	return b
}

func getAuth(username, password string) transport.AuthMethod {
	if username != "" || password != "" {
		return &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}
	return nil
}
