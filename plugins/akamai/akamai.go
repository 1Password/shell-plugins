package akamai

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AkamaiCLI() schema.Executable {
	return schema.Executable{
		Name:    "Akamai CLI",
		Runs:    []string{"akamai"},
		DocsURL: sdk.URL("https://techdocs.akamai.com/developer/docs/about-clis"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("config"),    // skip 1Password authentication for "akamai config" and its subcommands
			needsauth.NotWhenContainsArgs("install"),   // skip 1Password authentication for "akamai install" and its subcommands
			needsauth.NotWhenContainsArgs("get"),       // skip 1Password authentication for "akamai get" and its subcommands
			needsauth.NotWhenContainsArgs("list"),      // skip 1Password authentication for "akamai list" and its subcommands
			needsauth.NotWhenContainsArgs("search"),    // skip 1Password authentication for "akamai search" and its subcommands
			needsauth.NotWhenContainsArgs("uninstall"), // skip 1Password authentication for "akamai uninstall" and its subcommands
			needsauth.NotWhenContainsArgs("update"),    // skip 1Password authentication for "akamai update" and its subcommands
			needsauth.NotWhenContainsArgs("upgrade"),   // skip 1Password authentication for "akamai upgrade" and its subcommands
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIClientCredentials,
			},
		},
	}
}
