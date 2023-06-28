package render

import (
	"os"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)
func ConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "~/.render/config.yaml"
	}
	return configDir + "/render/config.yaml"
}

func RenderCLI() schema.Executable {
	return schema.Executable{
		Name:      "Render CLI", 
		Runs:      []string{"render"},
		DocsURL:   sdk.URL("https://render.com/docs/cli"), 
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
