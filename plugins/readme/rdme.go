package readme

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ReadMeCLI() schema.Executable {
	return schema.Executable{
		Name:    "ReadMe CLI",
		Runs:    []string{"rdme"},
		DocsURL: sdk.URL("https://docs.readme.com/main/docs/rdme"),
		NeedsAuth: needsauth.For(
			needsauth.NotForHelpOrVersion(),

			needsauth.NotWhenContainsArgs("--key"),

			needsauth.OnlyFor(
				needsauth.ForCommand("openapi"),

				needsauth.ForCommand("docs"),
				needsauth.ForCommand("docs:prune"),
				needsauth.ForCommand("guides"),
				needsauth.ForCommand("guides:prune"),

				needsauth.ForCommand("changelogs"),
				needsauth.ForCommand("custompages"),

				needsauth.ForCommand("versions"),
				needsauth.ForCommand("versions:create"),
				needsauth.ForCommand("versions:delete"),
				needsauth.ForCommand("versions:update"),

				needsauth.ForCommand("categories"),
				needsauth.ForCommand("categories:create"),
			),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
