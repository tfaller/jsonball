package startup

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/tfaller/jsonball"
	"github.com/tfaller/jsonball/internal/mysql"
)

const envRegistryConnectionString = "REGISTRY_CS"
const envRegistryDocumentKey = "REGISTRY_DOC_KEY"

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

	docKeyHex := MustGetEnvVar(envRegistryDocumentKey)
	docKey, err := hex.DecodeString(docKeyHex)
	if err != nil {
		return nil, fmt.Errorf("can't decode hex document key: %w", err)
	}

	return mysql.NewRegistry(cs, docKey)
}

// MustGetRegistry gets a registry, otherwise panics
func MustGetRegistry() jsonball.Registry {
	registry, err := GetRegistry()
	if err != nil {
		panic(err)
	}
	return registry
}
