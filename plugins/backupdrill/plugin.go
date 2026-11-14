package backupdrill

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "backupdrill",
		Platform: schema.PlatformInfo{
			Name:     "BackupDrill",
			Homepage: sdk.URL("https://backupdrill.com"),
		},
		Credentials: []schema.CredentialType{
			Credentials(),
		},
		Executables: []schema.Executable{
			BackupDrillCLI(),
		},
	}
}
