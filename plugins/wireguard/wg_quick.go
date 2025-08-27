package wireguard

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func WireguardVPNCLI() schema.Executable {
	return schema.Executable{
		Name:    "Wireguard VPN Quick CLI",
		Runs:    []string{"wg-quick"},
		DocsURL: sdk.URL("https://www.wireguard.com/quickstart/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessConfig,
			},
		},
	}
}
