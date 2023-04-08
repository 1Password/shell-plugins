package atlas

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "atlas",
		Platform: schema.PlatformInfo{
			Name:     "MongoDB Atlas",
			Homepage: sdk.URL("https://www.mongodb.com/"),
		},
		Credentials: []schema.CredentialType{
			PrivateKeyPair(),
		},
		Executables: []schema.Executable{
			MongoDBAtlasCLI(),
		},
	}
}
