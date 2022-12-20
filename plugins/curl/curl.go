package curl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/credselect"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func CurlCLI() schema.Executable {
	return schema.Executable{
		Name:      "cURL",
		Runs:      []string{"curl"},
		DocsURL:   sdk.URL("https://curl.se/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Select:        credselect.CanAuthenticateHTTPRequests,
				AllowMultiple: true,
				ProvisionerFunc: func(plugin string, credentialType schema.CredentialType) sdk.Provisioner {
					return CurlProvisioner(credentialType.HTTPProvisioner)
				},
			},
		},
	}
}
