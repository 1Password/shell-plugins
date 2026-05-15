package gcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func GsutilCLI() schema.Executable {
	return schema.Executable{
		Name:    "Google Cloud Storage CLI",
		Runs:    []string{"gsutil"},
		DocsURL: sdk.URL("https://cloud.google.com/storage/docs/gsutil"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.ServiceAccountKey,
			},
		},
	}
}
