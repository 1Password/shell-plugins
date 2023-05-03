package credname

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGettingCredentialIDsFromNames(t *testing.T) {
	names := []sdk.CredentialName{
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
	}

	expectedIDs := []string{
		"api_client_credentials",
		"api_key",
		"api_token",
		"access_key",
		"access_token",
		"app_password",
		"app_token",
		"auth_token",
		"cli_token",
		"credential",
		"credentials",
		"database_credentials",
		"login_details",
		"personal_api_token",
		"personal_access_token",
		"registry_credentials",
		"secret_key",
	}

	for i, name := range names {
		assert.Equal(t, expectedIDs[i], name.ID())
	}
}
