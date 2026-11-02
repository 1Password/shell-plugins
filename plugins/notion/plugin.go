package notion

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "notion",
		Platform: schema.PlatformInfo{
			Name:     "Notion",
			Homepage: sdk.URL("https://www.notion.com"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			NotionCLI(),
		},
	}
}
