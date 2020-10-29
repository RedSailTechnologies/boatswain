package kraken

import (
	"testing"

	pb "github.com/redsailtechnologies/boatswain/rpc/kraken"
	"github.com/stretchr/testify/assert"
)

func TestToHelmClient(t *testing.T) {
	sut := Cluster{
		&pb.Cluster{
			Name:     "cluster",
			Endpoint: "www.not.real",
			Token:    "abcdefg",
			Cert:     "notarealcert...",
		},
	}

	clientset, err := sut.ToHelmClient("")

	assert.NotNil(t, clientset)
	assert.Nil(t, err)
}

func TestToClientset(t *testing.T) {
	sut := Cluster{
		&pb.Cluster{
			Name:     "cluster",
			Endpoint: "www.not.real",
			Token:    "abcdefg",
			Cert:     "notarealcert...",
		},
	}

	clientset, err := sut.ToClientset()

	assert.NotNil(t, clientset)
	assert.Nil(t, err)
}
