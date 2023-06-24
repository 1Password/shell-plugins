package yugabytedb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func YugabyteDBCLI() schema.Executable {
	return schema.Executable{
		Name:      "YugabyteDB SQL Shell",
		Runs:      []string{"ysqlsh"},
		DocsURL:   sdk.URL("https://docs.yugabyte.com/preview/admin/ysqlsh/"), 
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DatabaseCredentials,
			},
		},
	}
}
