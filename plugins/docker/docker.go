package docker

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func DockerCLI() schema.Executable {
	return schema.Executable{
		Name:    "Docker CLI", // TODO: Check if this is correct
		Runs:    []string{"docker"},
		DocsURL: sdk.URL("https://docker.com/docs/cli"), // TODO: Replace with actual URL
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: "credname.UserLogin",
			},
		},
	}
}
