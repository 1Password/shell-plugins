package scorecard

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "scorecard",
		Platform: schema.PlatformInfo{
			Name:     "OpenSSF Scorecard",
			Homepage: sdk.URL("https://github.com/ossf/scorecard"),
		},
		Credentials: []schema.CredentialType{
			SecretKey(),
		},
		Executables: []schema.Executable{
			OpenSSFScorecardCLI(),
		},
	}
}
