package yugabytedb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "yugabytedb",
		Platform: schema.PlatformInfo{
			Name:     "YugabyteDB",
			Homepage: sdk.URL("https://yugabyte.com"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			YugabyteDBCLI(),
		},
	}
}
