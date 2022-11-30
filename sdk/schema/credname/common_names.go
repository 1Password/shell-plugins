package credname

import "github.com/1Password/shell-plugins/sdk"

// Credential type names.
const (
	AccessToken         = sdk.CredentialName("Access Token")
	APIKey              = sdk.CredentialName("API Key")
	APIToken            = sdk.CredentialName("API Token")
	PersonalAccessToken = sdk.CredentialName("Personal Access Token")
	PersonalAPIToken    = sdk.CredentialName("Personal API Token")
	CLIToken            = sdk.CredentialName("CLI Token")
	AuthToken           = sdk.CredentialName("Auth Token")
	AppToken            = sdk.CredentialName("App Token")
	AppPassword         = sdk.CredentialName("App Password")
	AccessKey           = sdk.CredentialName("Access Key")
	Credentials         = sdk.CredentialName("Credentials")
	DatabaseCredentials = sdk.CredentialName("Database Credentials")
	RegistryCredentials = sdk.CredentialName("Registry Credentials")
	Credential          = sdk.CredentialName("Credential")
	SecretKey           = sdk.CredentialName("Secret Key")
)
