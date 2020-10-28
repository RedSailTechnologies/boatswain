package kraken

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToHelmClient(t *testing.T) {
	config := ClusterConfig{
		Name:     "cluster",
		Endpoint: "www.not.real",
		Token:    "abcdefg",
		Cert:     "notarealcert...",
	}
	sut := Config{
		[]ClusterConfig{
			config,
		},
	}

	validClientset, noErr := sut.ToHelmClient("cluster", "")
	invalidClientset, err := sut.ToHelmClient("notreal", "")

	assert.NotNil(t, validClientset)
	assert.Nil(t, noErr)
	assert.Nil(t, invalidClientset)
	assert.Equal(t, errors.New("cluster not found"), err)
}

func TestToClientset(t *testing.T) {
	config := ClusterConfig{
		Name:     "cluster",
		Endpoint: "www.not.real",
		Token:    "abcdefg",
		Cert:     "notarealcert...",
	}
	sut := Config{
		[]ClusterConfig{
			config,
		},
	}

	validClientset, noErr := sut.ToClientset("cluster")
	invalidClientset, err := sut.ToClientset("notreal")

	assert.NotNil(t, validClientset)
	assert.Nil(t, noErr)
	assert.Nil(t, invalidClientset)
	assert.Equal(t, errors.New("cluster not found"), err)
}

func TestGetClusterConfig(t *testing.T) {
	config := ClusterConfig{
		Name:     "cluster",
		Endpoint: "www.not.real",
		Token:    "abcdefg",
		Cert:     "notarealcert...",
	}
	sut := Config{
		[]ClusterConfig{
			config,
		},
	}

	valid, noErr := sut.getClusterConfig("cluster")
	invalid, err := sut.getClusterConfig("doesn'texist")

	assert.Equal(t, &config, valid)
	assert.Nil(t, noErr)
	assert.Nil(t, invalid)
	assert.Equal(t, errors.New("cluster not found"), err)
}
