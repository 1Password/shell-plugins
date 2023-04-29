package snowflake

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "snowflake",
		Platform: schema.PlatformInfo{
			Name:     "Snowflake",
			Homepage: sdk.URL("https://snowflake.com"),
		},
		Credentials: []schema.CredentialType{
			LoginDetails(),
		},
		Executables: []schema.Executable{
			SnowflakeCLI(),
		},
	}
}
