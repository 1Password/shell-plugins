package sqitch

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func Sqitch() schema.Executable {
	return schema.Executable{
		Name:      "Sqitch",
		Runs:      []string{"sqitch"},
		DocsURL:   sdk.URL("https://sqitch.org/docs/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.IfAny(
				needsauth.ForCommand("deploy"),
				needsauth.ForCommand("log"),
				needsauth.ForCommand("revert"),
				needsauth.ForCommand("status"),
				needsauth.ForCommand("upgrade"),
				needsauth.ForCommand("verify"),
			),
			needsauth.NotForHelp(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DatabaseCredentials,
			},
		},
	}
}
