package template

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/redsailtechnologies/boatswain/rpc/repo"
)

type MockRepo struct{}

func (r *MockRepo) Create(context.Context, *repo.CreateRepo) (*repo.RepoCreated, error) {
	return nil, nil
}

func (r *MockRepo) Update(context.Context, *repo.UpdateRepo) (*repo.RepoUpdated, error) {
	return nil, nil
}

func (r *MockRepo) Destroy(context.Context, *repo.DestroyRepo) (*repo.RepoDestroyed, error) {
	return nil, nil
}

func (r *MockRepo) Read(context.Context, *repo.ReadRepo) (*repo.RepoRead, error) {
	return nil, nil
}

func (r *MockRepo) Find(context.Context, *repo.FindRepo) (*repo.RepoFound, error) {
	return &repo.RepoFound{
		Uuid: "1",
	}, nil
}

func (r *MockRepo) All(context.Context, *repo.ReadRepos) (*repo.ReposRead, error) {
	return nil, nil
}

func (r *MockRepo) Chart(context.Context, *repo.ReadChart) (*repo.ChartRead, error) {
	return nil, nil
}

func (r *MockRepo) File(context.Context, *repo.ReadFile) (*repo.FileRead, error) {
	return &repo.FileRead{
		File: []byte(`
name: sample-app-two
helm:
  chart: sample-app-two
  repo: another-repo
  version: ${{ .Inputs.version }} 
`),
	}, nil
}

func TestUnmarshal(t *testing.T) {
	p, err := filepath.Abs("../../../docs/example.yaml")
	f, err := ioutil.ReadFile(p)
	sut := NewEngine(context.TODO(), &MockRepo{})
	if err == nil {
		sut.Run(f, []byte(`
sampleAppChartVersion: '1.1.1'
someValue: somereplacedvalue
version: '0.1.0'
`))
	}
}
