package zapier

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ZapierCLI() schema.Executable {
	return zapierExecutable("zapier")
}

// ZapierPlatformCLI covers zapier-platform-cli v19+, which no longer installs a `zapier` binary.
func ZapierPlatformCLI() schema.Executable {
	return zapierExecutable("zapier-platform")
}

func zapierExecutable(command string) schema.Executable {
	return schema.Executable{
		Name:    "Zapier CLI",
		Runs:    []string{command},
		DocsURL: sdk.URL("https://platform.zapier.com/cli_docs/docs"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DeployKey,
			},
		},
	}
}
