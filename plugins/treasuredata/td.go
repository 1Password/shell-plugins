package treasuredata

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func TreasureDataCLI() schema.Executable {
	return schema.Executable{
		Name:    "Treasure Data Toolbelkt",
		Runs:    []string{"td"},
		DocsURL: sdk.URL("https://docs.treasuredata.com/display/public/PD/TD+Toolbelt"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessKey,
			},
		},
	}
}
