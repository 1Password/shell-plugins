package postman

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "postman",
		Platform: schema.PlatformInfo{
			Name:     "postman",
			Homepage: sdk.URL("https://www.postman.com/"), 
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			postmanCLI(),
		},
	}
}
