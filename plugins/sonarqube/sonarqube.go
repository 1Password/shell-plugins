package sonarqube

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SonarQubeCLI() schema.Executable {
	return schema.Executable{
		Name:    "SonarQube CLI",
		Runs:    []string{"sonar"},
		DocsURL: sdk.URL("https://docs.sonarsource.com/sonarqube-cli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForExactArgs("auth", "login"),
			needsauth.NotForExactArgs("auth", "logout"),
			needsauth.NotWhenContainsArgs("config", "telemetry"),
			needsauth.NotWhenContainsArgs("system"),
			needsauth.NotWhenContainsArgs("self-update"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.CLIToken,
			},
		},
	}
}
