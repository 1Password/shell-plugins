package b2

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func BackblazeB2CLI() schema.Executable {
	return schema.Executable{
		Name:    "Backblaze B2 Cloud Storage Command-Line Tools",
		Runs:    []string{"b2"},
		DocsURL: sdk.URL("https://www.backblaze.com/docs/cloud-storage-command-line-interface"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.ApplicationKey,
			},
		},
	}
}
