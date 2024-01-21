package npm

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func NPMCLI() schema.Executable {
	return schema.Executable{
		Name:    "NPM CLI",
		Runs:    []string{"npm"},
		DocsURL: sdk.URL("https://docs.npmjs.com/cli"),
		NeedsAuth: needsauth.IfAll(
			// not a complete list of commands that don't
			// need auth, but probably the main ones
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWhenContainsArgs("init"),
			needsauth.NotWhenContainsArgs("?"),
			needsauth.NotWhenContainsArgs("config"),
			needsauth.NotWhenContainsArgs("help-search"),
			needsauth.NotWhenContainsArgs("login"),
			needsauth.NotWhenContainsArgs("logout"),
			needsauth.NotWhenContainsArgs("prune"),
			needsauth.NotWhenContainsArgs("shrinkwrap"),
			needsauth.NotWhenContainsArgs("start"),
			needsauth.NotWhenContainsArgs("run-script"),
			needsauth.NotWhenContainsArgs("run"),
			needsauth.NotWhenContainsArgs("rum"),
			needsauth.NotWhenContainsArgs("urn"),
			needsauth.NotWhenContainsArgs("uninstall"),
			needsauth.NotWhenContainsArgs("unlink"),
			needsauth.NotWhenContainsArgs("remove"),
			needsauth.NotWhenContainsArgs("rm"),
			needsauth.NotWhenContainsArgs("r"),
			needsauth.NotWhenContainsArgs("un"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
