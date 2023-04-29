package treasuredata

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func TreasureDataCLI() schema.Executable {
	return schema.Executable{
		Name:    "Treasure Data Toolbelt",
		Runs:    []string{"td"},
		DocsURL: sdk.URL("https://docs.treasuredata.com/display/public/PD/TD+Toolbelt"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("-k"),
			needsauth.NotWhenContainsArgs("-c"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
