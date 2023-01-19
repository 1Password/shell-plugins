package databricks

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "databricks",
		Platform: schema.PlatformInfo{
			Name:     "Databricks",
			Homepage: sdk.URL("https://databricks.com"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			DatabricksCLI(),
		},
	}
}
