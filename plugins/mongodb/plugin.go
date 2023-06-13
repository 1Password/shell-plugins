package mongodb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "mongodb",
		Platform: schema.PlatformInfo{
			Name:     "MongoDB Shell",
			Homepage: sdk.URL("https://www.mongodb.com/products/shell"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			MongoshCLI(),
		},
	}
}
