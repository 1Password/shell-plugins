package credselect

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

const (
	Any                             = sdk.CredentialSelector("any")
	SAMLIdentityProvider            = sdk.CredentialSelector("saml")
	CanAuthenticateHTTPRequests     = sdk.CredentialSelector("http")
	CanAuthenticateToDockerRegistry = sdk.CredentialSelector("docker-registry")
)

var details = map[sdk.CredentialSelector]CredentialSelectorDetails{
	Any: {
		Description: "Any credential from any plugin",
		Matches: func(credentialType schema.CredentialType) bool {
			return true
		},
	},
	SAMLIdentityProvider: {
		Description: "Credentials that can be used for SAML authentication",
	},
	CanAuthenticateHTTPRequests: {
		Description: "Credentials that can authenticate HTTP requests",
		Matches: func(credentialType schema.CredentialType) bool {
			return credentialType.HTTPProvisioner != nil
		},
	},
	CanAuthenticateToDockerRegistry: {
		Description: "Credentials that can be used to authenticate to a Docker Registry",
	},
}

type CredentialSelectorDetails struct {
	Description string
	Matches     func(credentialType schema.CredentialType) bool
}
