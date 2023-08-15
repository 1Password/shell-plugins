package todoist

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "todoist",
		Platform: schema.PlatformInfo{
			Name:     "Todoist",
			Homepage: sdk.URL("https://todoist.com"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			TodoistCLI(),
		},
	}
}
