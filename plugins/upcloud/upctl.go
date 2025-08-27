package upcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func UpCloudCLI() schema.Executable {
	return schema.Executable{
		Name:    "UpCloud CLI",
		Runs:    []string{"upctl"},
		DocsURL: sdk.URL("https://upcloudltd.github.io/upcloud-cli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("completion"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.UserLogin,
			},
		},
	}
}
