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
			Homepage: sdk.URL("https://www.mongodb.com/"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			Credentials(),
		},
		Executables: []schema.Executable{
			MongoDBAtlasCLI(),
		},
	}
}
