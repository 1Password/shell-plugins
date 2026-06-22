package pypi

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "pypi",
		Platform: schema.PlatformInfo{
			Name:     "PyPI",
			Homepage: sdk.URL("https://pypi.org"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			TwineCLI(),
			FlitCLI(),
			HatchCLI(),
		},
	}
}
