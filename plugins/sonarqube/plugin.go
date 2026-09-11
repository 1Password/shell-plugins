package sonarqube

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "sonarqube",
		Platform: schema.PlatformInfo{
			Name:     "SonarQube CLI",
			Homepage: sdk.URL("https://docs.sonarsource.com/sonarqube-cli/"),
		},
		Credentials: []schema.CredentialType{
			CLIToken(),
		},
		Executables: []schema.Executable{
			SonarQubeCLI(),
		},
	}
}
