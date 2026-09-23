package cline

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "cline",
		Platform: schema.PlatformInfo{
			Name:     "Cline",
			Homepage: sdk.URL("https://cline.bot"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			ClineCLI(),
		},
	}
}
