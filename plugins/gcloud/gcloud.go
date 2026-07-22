package gcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func GCloudCLI() schema.Executable {
	return schema.Executable{
		Name:    "Google Cloud CLI",
		Runs:    []string{"gcloud"},
		DocsURL: sdk.URL("https://cloud.google.com/sdk/gcloud"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("auth"),
			needsauth.NotWhenContainsArgs("config"),
			needsauth.NotWhenContainsArgs("info"),
			needsauth.NotWhenContainsArgs("components"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.ServiceAccountKey,
			},
		},
	}
}
