package kaggle

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "kaggle",
		Platform: schema.PlatformInfo{
			Name:     "Kaggle",
			Homepage: sdk.URL("https://kaggle.com"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			KaggleCLI(),
		},
	}
}
