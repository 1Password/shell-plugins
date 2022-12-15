package homebrew

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

var commands = [][]string{
	{"search"},

	{"bump-cask-pr"},
	{"bump-formula-pr"},
}

func HomebrewCLI() schema.Executable {
	return schema.Executable{
		Name:      "Homebrew CLI",
		Runs:      []string{"brew"},
		DocsURL:   sdk.URL("https://brew.sh/"),
		NeedsAuth: needsauth.ForCommands(commands...),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
