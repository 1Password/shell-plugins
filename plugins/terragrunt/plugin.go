package terragrunt

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

// This plugin is based on the Terraform plugin.

func New() schema.Plugin {
	return schema.Plugin{
		Name: "terragrunt",
		Platform: schema.PlatformInfo{
			Name:     "Terragrunt",
			Homepage: sdk.URL("https://terragrunt.gruntwork.io"),
		},
		Executables: []schema.Executable{
			TerragruntCLI(),
		},
	}
}
