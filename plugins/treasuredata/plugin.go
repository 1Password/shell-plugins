package treasuredata

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "treasuredata",
		Platform: schema.PlatformInfo{
			Name:     "Treasure Data",
			Homepage: sdk.URL("https://www.treasuredata.com/"),
		},
		Credentials: []schema.CredentialType{
			AccessKey(),
		},
		Executables: []schema.Executable{
			TreasureDataCLI(),
		},
	}
}
