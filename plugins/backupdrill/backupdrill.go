package backupdrill

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func BackupDrillCLI() schema.Executable {
	return schema.Executable{
		Name:    "BackupDrill CLI",
		Runs:    []string{"backupdrill"},
		DocsURL: sdk.URL("https://backupdrill.com/docs/cli"),
		// Only treat `help` as the help command when it comes first. NotForHelp() would also skip
		// auth for a value like `backup --project-name help`.
		NeedsAuth: needsauth.IfAll(
			needsauth.NotWhenContainsArgs("-h"),
			needsauth.NotWhenContainsArgs("--help"),
			needsauth.NotForCommand("help"),
			needsauth.NotForVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credentials,
			},
		},
	}
}
