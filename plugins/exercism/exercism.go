package exercism

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ExercismCLI() schema.Executable {
	return schema.Executable{
		Name:    "Exercism CLI",
		Runs:    []string{"exercism"},
		DocsURL: sdk.URL("https://exercism.org/docs/using/solving-exercises/working-locally"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForCommand("completion"),
			needsauth.NotForCommand("upgrade"),
			needsauth.NotForCommand("workspace"),
			needsauth.NotForCommand("troubleshoot"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
