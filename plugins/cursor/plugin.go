package cursor

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "cursor",
		Platform: schema.PlatformInfo{
			Name:     "Cursor",
			Homepage: sdk.URL("https://cursor.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			CursorCLI(),
		},
	}
}
