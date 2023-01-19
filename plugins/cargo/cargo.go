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
		NeedsAuth: needsauth.For(
			needsauth.OnlyFor(
				needsauth.ForCommand("publish"),
				needsauth.ForCommand("yank"),
				needsauth.ForCommand("owner"),
			),
			needsauth.NotForHelp(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
