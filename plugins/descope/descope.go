package descope

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func DesopeCLI() schema.Executable {
	return schema.Executable{
		Name:      "Desope CLI", // TODO: Check if this is correct
		Runs:      []string{"descope"},
		DocsURL:   sdk.URL("https://descope.com/docs/cli"), // TODO: Replace with actual URL
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.&#34;ManagementKey&#34;,
			},
		},
	}
}
