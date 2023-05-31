package credname

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/stretchr/testify/assert"
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
		Password,
		PersonalAPIToken,
		PersonalAccessToken,
		RegistryCredentials,
		SecretKey,
		sdk.CredentialName(""),
		sdk.CredentialName("Database-specific Credentials"),
		sdk.CredentialName("Public/Private Key-Pair"),
		sdk.CredentialName("Public/Private Key Pair"),
		sdk.CredentialName("KeyPair"),
		sdk.CredentialName("___some-test/name "),
		sdk.CredentialName("this -is/NOT-a, real. _- /  life---- scENARIo..."),
		sdk.CredentialName("___some,.other test/name-which -is/NOT-a, real. , - ./ - life_scENARIo..."),
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
		"password",
		"personal_api_token",
		"personal_access_token",
		"registry_credentials",
		"secret_key",
		"",
		"database_specific_credentials",
		"public_private_key_pair",
		"public_private_key_pair",
		"keypair",
		"some_test_name",
		"this_is_not_a_real_life_scenario",
		"some_other_test_name_which_is_not_a_real_life_scenario",
	}

	for i, name := range names {
		assert.Equal(t, expectedIDs[i], name.ID().String())
	}
}
