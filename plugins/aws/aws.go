package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/credselect"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AWSCLI() schema.Executable {
	return schema.Executable{
		Name:      "AWS CLI",
		Runs:      []string{"aws"},
		DocsURL:   sdk.URL("https://aws.amazon.com/cli/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				SelectFrom: []schema.CredentialUsage{
					{
						Name: credname.AccessKey,
					},
					{
						Select: credselect.SAMLIdentityProvider,
						ProvisionerFunc: func(plugin string, credentialType schema.CredentialType) sdk.Provisioner {
							// TODO: Return SAML AWS provisioner
							return nil
						},
					},
				},
			},
		},
	}
}
