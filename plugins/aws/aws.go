package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func Executable_aws() schema.Executable {
	return schema.Executable{
		Runs:      []string{"aws"},
		Name:      "AWS CLI",
		DocsURL:   sdk.URL("https://aws.amazon.com/cli/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Credentials: []schema.CredentialType{
			AccessKey(),
		},
	}
}
