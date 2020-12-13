package startup

import (
	"fmt"
	"os"
)

// ErrEnvVarNotSet indicates that an important environment variable is missing
type ErrEnvVarNotSet string

func (e ErrEnvVarNotSet) Error() string {
	return fmt.Sprintf("environment variable %q is not set", string(e))
}

// MustGetEnvVar gets an environment variable, otherwise panics
func MustGetEnvVar(name string) string {
	val := os.Getenv(name)
	if val == "" {
		panic(ErrEnvVarNotSet(name))
	}
	return val
}
