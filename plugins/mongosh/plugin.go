package mongosh

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "mongosh",
		Platform: schema.PlatformInfo{
			Name:     "MongoDB Shell",
			Homepage: sdk.URL("https://mongosh.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			MongoDBShellCLI(),
		},
	}
}
