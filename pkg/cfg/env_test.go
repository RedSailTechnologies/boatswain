package cfg

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvOrDefaultString(t *testing.T) {
	key := "someenvvar"
	def1 := "default1"
	val1 := "val1"
	os.Setenv(key, val1)
	assert.Equal(t, val1, EnvOrDefaultString(key, def1))
	assert.Equal(t, "a", EnvOrDefaultString("ne", "a"))
}

func TestEnvOrDefaultInt(t *testing.T) {
	key := "someenvvar"
	def1 := 1
	val1 := 2
	os.Setenv(key, fmt.Sprint(val1))
	assert.Equal(t, val1, EnvOrDefaultInt(key, def1))
	assert.Equal(t, 3, EnvOrDefaultInt("ne", 3))
}
