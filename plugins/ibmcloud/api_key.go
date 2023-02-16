package ibmcloud

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/iam"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/config_helpers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		DocsURL:       sdk.URL("https://cloud.ibm.com/docs/account?topic=account-userapikey"),
		ManagementURL: sdk.URL("https://cloud.ibm.com/iam/apikeys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to IBM Cloud.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 44,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.TempFile(ibmCloudConfig, provision.Filename("config.json")),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"IBMCLOUD_API_KEY": fieldname.APIKey,
}

func ibmCloudConfig(in sdk.ProvisionInput) ([]byte, error) {
	// Config info to cache for more efficient authentication
	// These values get added and removed from the config file when calling `ibmcloud login` and `ibmcloud logout`
	type ibmCacheEntry struct {
		accessToken    string
		refreshToken   string
		accountDetails models.Account
	}

	var cache ibmCacheEntry
	// Try finding the config info in cache
	ok := in.Cache.Get("IBMAccessToken", cache.accessToken) &&
		in.Cache.Get("IBMRefreshToken", cache.refreshToken) &&
		in.Cache.Get("IBMAccountDetails", cache.accountDetails)
	if !ok {
		// If the config info is not in cache, retrieve it using the IBM Cloud SDK
		tokenRequest := iam.APIKeyTokenRequest(in.ItemFields[fieldname.APIKey])
		restClient := iam.NewClient(iam.DefaultConfig("https://iam.cloud.ibm.com"), rest.NewClient()) // TODO: Check whether the iamEndpoint should be fetched instead of hard-coded
		token, err := restClient.GetToken(tokenRequest)
		if err != nil {
			return nil, err
		}
		cache.accessToken = token.AccessToken
		cache.refreshToken = token.RefreshToken
		// TODO: Is there a way to retrieve the account information using the API key?
		// cache.accountDetails.GUID = ""
		// cache.accountDetails.Name = ""
		// cache.accountDetails.Owner = ""
	}

	// Get the config file path
	// If the `IBM_CLOUD_HOME` env var is set use that, otherwise default to the user's home directory
	var configDir string
	if configDir = bluemix.EnvConfigHome.Get(); configDir == "" {
		var homeDir string
		if homeDir = config_helpers.UserHomeDir(); homeDir == "" {
			return nil, errors.New("could not retrieve user home directory")
		}
		configDir = filepath.Join(homeDir, "/.bluemix")
	}
	configPath := filepath.Join(configDir, "config.json")

	var initialConfig core_config.BXConfigData

	// Check if an already existing config file is present and if so use its information to build the temporary config file
	if _, err := os.Stat(configPath); err == nil {
		configFile, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(configFile, &initialConfig)
		if err != nil {
			return nil, err
		}
	} else {
		// TODO: Reconsider whether these should be hard-coded
		initialConfig.APIEndpoint = "https://cloud.ibm.com"
		initialConfig.ConsoleEndpoint = "https://cloud.ibm.com"
		initialConfig.CloudType = "public"
		initialConfig.CloudName = "bluemix"
		initialConfig.IAMEndpoint = "https://iam.cloud.ibm.com"
	}

	initialConfig.IAMToken = "Bearer " + cache.accessToken
	initialConfig.IAMRefreshToken = cache.refreshToken
	initialConfig.Account.GUID = cache.accountDetails.GUID
	initialConfig.Account.Name = cache.accountDetails.Name
	initialConfig.Account.Owner = cache.accountDetails.Owner

	configJSON, err := json.Marshal(initialConfig)
	if err != nil {
		return nil, err
	}

	// Uncomment below for a quick way to check the JSON output
	// return nil, fmt.Errorf(string(configJSON))

	// TODO: Write config info to cache

	return configJSON, nil
}
