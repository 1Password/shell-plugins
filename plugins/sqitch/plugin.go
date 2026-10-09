package sqitch

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "sqitch",
		Platform: schema.PlatformInfo{
			Name:     "Sqitch",
			Homepage: sdk.URL("https://sqitch.org"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			Sqitch(),
		},
	}
}
