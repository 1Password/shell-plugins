package docker

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/credselect"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func DockerCLI() schema.Executable {
	return schema.Executable{
		Name:      "Docker CLI",
		Runs:      []string{"docker"},
		DocsURL:   sdk.URL("https://docs.docker.com/engine/reference/commandline/cli/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Description: "Authenticate to a Docker Registry for the push and pull commands.",
				Select:      credselect.CanAuthenticateToDockerRegistry,
				NeedsAuth: needsauth.ForCommands(
					[]string{"pull"},
					[]string{"push"},
				),
			},
			{
				Description: "Provision run command with secrets.",
				Select:      credselect.Any,
				NeedsAuth: needsauth.ForCommands(
					[]string{"run"},
				),
				AllowMultiple: true,
				ProvisionerFunc: func(plugin string, credentialType schema.CredentialType) sdk.Provisioner {
					return RunCommandProvisioner(credentialType.DefaultProvisioner)
				},
			},
		},
	}
}
