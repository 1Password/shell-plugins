package terraform

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "terraform",
		Platform: schema.PlatformInfo{
			Name:     "Terraform",
			Homepage: sdk.URL("https://terraform.io"),
		},
		Executables: []schema.Executable{
			TerraformCLI(),
		},
	}
}
