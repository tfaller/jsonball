package startup

import (
	"fmt"
	"os"

	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/internal/mysql"
)

const envRegistryConnectionString = "REGISTRY_CS"

// GetRegistry gets a registry
func GetRegistry() (jsonball.Registry, error) {
	cs := os.Getenv(envRegistryConnectionString)
	if cs == "" {
		return nil, fmt.Errorf(
			"Expected env %q to have format %q",
			envRegistryConnectionString,
			"[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]",
		)
	}
	return mysql.NewRegistry(cs)
}

// MustGetRegistry gets a registry, otherwise panics
func MustGetRegistry() jsonball.Registry {
	registry, err := GetRegistry()
	if err != nil {
		panic(err)
	}
	return registry
}
