package cargo

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

var commands = [][]string{
	{"publish"},
	{"yank"},
	{"owner"},
}

func CargoCLI() schema.Executable {
	return schema.Executable{
		Name:      "Cargo CLI",
		Runs:      []string{"cargo"},
		DocsURL:   sdk.URL("https://doc.rust-lang.org/cargo/index.html"),
		NeedsAuth: needsauth.ForCommands(commands...),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
