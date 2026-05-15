package plugins

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/1Password/shell-plugins/sdk/schema"
)

var registry = []schema.Plugin{}

func RegistryJSON() ([]byte, error) {
	var listToJSON []schema.Plugin
	// Only include plugins that contain at least one credential.
	for _, p := range List() {
		if len(p.Credentials) > 0 {
			listToJSON = append(listToJSON, p)
		}
	}
	return json.MarshalIndent(listToJSON, "", "\t")
}

func List() []schema.Plugin {
	list := make([]schema.Plugin, len(registry))
	copy(list, registry)
	return list
}

func Get(pluginName string) (schema.Plugin, error) {
	for _, p := range registry {
		if p.Name == pluginName {
			return p, nil
		}
	}
	return schema.Plugin{}, fmt.Errorf("unknown plugin: %s", pluginName)
}

func GetByExecutable(executableQuery string) (schema.Plugin, schema.Executable, error) {
	for _, p := range registry {
		for _, e := range p.Executables {
			if strings.EqualFold(executableQuery, e.Command()) || strings.EqualFold(executableQuery, e.Name) {
				return p, e, nil
			}
		}
	}
	return schema.Plugin{}, schema.Executable{}, fmt.Errorf("unknown plugin: %s", executableQuery)
}

func GetCredentialType(pluginName string, credentialName string) (schema.CredentialType, error) {
	for _, p := range registry {
		if p.Name == pluginName {
			for _, credential := range p.Credentials {
				return credential, nil
			}
			return schema.CredentialType{}, fmt.Errorf("unknown credential: %s (%s)", credentialName, pluginName)
		}
	}
	return schema.CredentialType{}, fmt.Errorf("unknown plugin: %s", pluginName)
}

func Register(p schema.Plugin) {
	registry = append(registry, p)
}
