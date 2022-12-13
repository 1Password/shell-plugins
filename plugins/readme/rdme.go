package readme

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

var commands = [][]string{
	[]string{"openapi"},

	[]string{"docs"},
	[]string{"docs:prune"},
	[]string{"guides"},
	[]string{"guides:prune"},

	[]string{"changelogs"},
	[]string{"custompages"},

	[]string{"versions"},
	[]string{"versions:create"},
	[]string{"versions:delete"},
	[]string{"versions:update"},

	[]string{"categories"},
	[]string{"categories:create"},
}

func ReadMeCLI() schema.Executable {
	return schema.Executable{
		Name:      "ReadMe CLI",
		Runs:      []string{"rdme"},
		DocsURL:   sdk.URL("https://docs.readme.com/main/docs/rdme"),
		NeedsAuth: needsauth.ForCommands(commands),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
