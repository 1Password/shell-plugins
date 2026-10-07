package cockroachdb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "cockroachdb",
		Platform: schema.PlatformInfo{
			Name:     "CockroachDB",
			Homepage: sdk.URL("https://www.cockroachlabs.com"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			Cockroach(),
		},
	}
}
