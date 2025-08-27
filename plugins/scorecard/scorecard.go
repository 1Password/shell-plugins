package scorecard

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OpenSSFScorecardCLI() schema.Executable {
	return schema.Executable{
		Name:    "OpenSSF Scorecard CLI",
		Runs:    []string{"scorecard"},
		DocsURL: sdk.URL("https://github.com/ossf/scorecard"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.SecretKey,
			},
		},
	}
}
