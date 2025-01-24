package motherduck

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "motherduck",
		Platform: schema.PlatformInfo{
			Name:     "MotherDuck",
			Homepage: sdk.URL("https://motherduck.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			DuckDBCLI(),
		},
	}
}
