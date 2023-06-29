package firebase

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "firebase",
		Platform: schema.PlatformInfo{
			Name:     "firebase",
			Homepage: sdk.URL("https://firebase.google.com/"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			firebaseCLI(),
		},
	}
}
