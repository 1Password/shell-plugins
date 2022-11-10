package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AWSCLI() schema.Executable {
	return schema.Executable{
		Runs:      []string{"aws"},
		Name:      "AWS CLI",
		DocsURL:   sdk.URL("https://aws.amazon.com/cli/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		UsesCredentials: []schema.CredentialUsage{
			{
				Name: credname.AccessKey,
			},
		},
	}
}
