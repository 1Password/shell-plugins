package cargo

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func CargoCLI() schema.Executable {
	return schema.Executable{
		Name:    "Cargo CLI",
		Runs:    []string{"cargo"},
		DocsURL: sdk.URL("https://doc.rust-lang.org/cargo/index.html"),
		NeedsAuth: needsauth.OnlyFor(
			needsauth.ForCommand("publish"),
			needsauth.ForCommand("yank"),
			needsauth.ForCommand("owner"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
