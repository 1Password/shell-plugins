package gcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func BqCLI() schema.Executable {
	return schema.Executable{
		Name:    "BigQuery CLI",
		Runs:    []string{"bq"},
		DocsURL: sdk.URL("https://cloud.google.com/bigquery/docs/bq-command-line-tool"),
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
