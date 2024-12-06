package ssh

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SSHCLI() schema.Executable {
	return schema.Executable{
		Name:    "SSH CLI",
		Runs:    []string{"sshpass"},
		DocsURL: sdk.URL("https://linux.die.net/man/1/sshpass"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.UserLogin,
			},
		},
	}
}
