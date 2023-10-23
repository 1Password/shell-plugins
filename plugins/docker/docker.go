package docker

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func DockerCLI() schema.Executable {
	return schema.Executable{
		Name:    "Docker CLI",
		Runs:    []string{"docker"},
		DocsURL: sdk.URL("https://docs.docker.com/engine/reference/commandline/docker/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			// NOTE for some reason, doesn't work when this isn't commented out
			//needsauth.ForCommand("login", "push", "pull"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.Credentials,
				Provisioner: dockerProvisioner{},
			},
		},
	}
}
