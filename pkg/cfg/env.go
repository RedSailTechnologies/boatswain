package cfg

import (
	"os"
	"strconv"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
)

// EnvOrDefaultString takes an env var and returns its value as a string or the default if its not set
func EnvOrDefaultString(e, d string) string {
	if val, ok := os.LookupEnv(e); ok {
		return val
	}
	return d
}

// EnvOrDefaultInt takes an env var and returns its value as an int or the default if its not set
func EnvOrDefaultInt(e string, d int) int {
	if val, ok := os.LookupEnv(e); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			logger.Warn("error converting int value", "error", err)
		}
		return v
	}
	return d
}
