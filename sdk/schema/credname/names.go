package credname

import "github.com/1Password/shell-plugins/sdk"

// Credential type names.
const (
	APIKey              = sdk.CredentialName("API Key")
	APIToken            = sdk.CredentialName("API Token")
	AccessKey           = sdk.CredentialName("Access Key")
	AccessToken         = sdk.CredentialName("Access Token")
	AppPassword         = sdk.CredentialName("App Password")
	AppToken            = sdk.CredentialName("App Token")
	AuthToken           = sdk.CredentialName("Auth Token")
	CLIToken            = sdk.CredentialName("CLI Token")
	Credential          = sdk.CredentialName("Credential")
	Credentials         = sdk.CredentialName("Credentials")
	DatabaseCredentials = sdk.CredentialName("Database Credentials")
	PersonalAPIToken    = sdk.CredentialName("Personal API Token")
	PersonalAccessToken = sdk.CredentialName("Personal Access Token")
	RegistryCredentials = sdk.CredentialName("Registry Credentials")
	SecretKey           = sdk.CredentialName("Secret Key")
)
