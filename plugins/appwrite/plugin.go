package appwrite

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "appwrite",
		Platform: schema.PlatformInfo{
			Name:     "Appwrite",
			Homepage: sdk.URL("https://appwrite.io"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			AppwriteCLI(),
		},
	}
}
