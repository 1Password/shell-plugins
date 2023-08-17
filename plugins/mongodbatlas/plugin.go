package mongodbatlas

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "mongodbatlas",
		Platform: schema.PlatformInfo{
			Name:     "MongoDB Atlas",
			Homepage: sdk.URL("https://www.mongodb.com/"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			MongoDBAtlasCLI(),
		},
	}
}
