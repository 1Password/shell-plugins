package docker

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "docker",
		Platform: schema.PlatformInfo{
			Name:     "Docker",
			Homepage: sdk.URL("https://www.docker.com/"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			DockerCLI(),
		},
	}
}
