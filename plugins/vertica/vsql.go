package vertica

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func VerticaCLI() schema.Executable {
	return schema.Executable{
		Name:    "Vertica CLI",
		Runs:    []string{"vsql"},
		DocsURL: sdk.URL("https://www.vertica.com/docs/9.2.x/HTML/Content/Authoring/ConnectingToVertica/vsql/UsingVsql.htm"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DatabaseCredentials,
			},
		},
	}
}
