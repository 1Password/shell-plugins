package openbao

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "openbao",
		Platform: schema.PlatformInfo{
			Name:     "OpenBao",
			Homepage: sdk.URL("https://openbao.org"),
		},
		Credentials: []schema.CredentialType{
			AuthToken(),
		},
		Executables: []schema.Executable{
			OpenBaoCLI(),
		},
	}
}
