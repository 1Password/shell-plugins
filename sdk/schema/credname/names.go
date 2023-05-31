package credname

import "github.com/1Password/shell-plugins/sdk"

// Credential type names.
const (
	APIClientCredentials = sdk.CredentialName("API Client Credentials")
	APIKey               = sdk.CredentialName("API Key")
	APIToken             = sdk.CredentialName("API Token")
	AccessKey            = sdk.CredentialName("Access Key")
	AccessToken          = sdk.CredentialName("Access Token")
	AppPassword          = sdk.CredentialName("App Password")
	AppToken             = sdk.CredentialName("App Token")
	AuthToken            = sdk.CredentialName("Auth Token")
	CLIToken             = sdk.CredentialName("CLI Token")
	Credential           = sdk.CredentialName("Credential")
	Credentials          = sdk.CredentialName("Credentials")
	DatabaseCredentials  = sdk.CredentialName("Database Credentials")
	LoginDetails         = sdk.CredentialName("Login Details")
	PersonalAPIToken     = sdk.CredentialName("Personal API Token")
	PersonalAccessToken  = sdk.CredentialName("Personal Access Token")
	RegistryCredentials  = sdk.CredentialName("Registry Credentials")
	SecretKey            = sdk.CredentialName("Secret Key")
	UserPass             = sdk.CredentialName("Username and Password")
)

func ListAll() []sdk.CredentialName {
	return []sdk.CredentialName{
		APIClientCredentials,
		APIKey,
		APIToken,
		AccessKey,
		AccessToken,
		AppPassword,
		AppToken,
		AuthToken,
		CLIToken,
		Credential,
		Credentials,
		DatabaseCredentials,
		LoginDetails,
		PersonalAPIToken,
		PersonalAccessToken,
		RegistryCredentials,
		SecretKey,
		UserPass,
	}
}
