package databricks

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func DatabricksCLI() schema.Executable {
	return schema.Executable{
		Name:      "Databricks CLI",
		Runs:      []string{"databricks"},
		DocsURL:   sdk.URL("https://docs.databricks.com/dev-tools/cli/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
