package ibmcloud

import (
	"context"
	"encoding/json"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/config_helpers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"os"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/uaa"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
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
		DefaultProvisioner: provision.TempFile(ibmCloudConfig, provision.Filename("config.json"), provision.SetOutputDirAsEnvVar("IBMCLOUD_HOME")),
		Importer: importer.TryAll(
			TryIBMCloudConfigFile(),
		)}
}

func ibmCloudConfig(in sdk.ProvisionInput) ([]byte, error) {
	var accessToken string
	// try finding it in cache
	ok := in.Cache.Get("IBMAccessToken", accessToken)
	if !ok {
		// if the credential is not in cache, fetch it from ibm
		tokenRequest := uaa.APIKeyTokenRequest(in.ItemFields[fieldname.APIKey])
		restClient := uaa.NewClient(uaa.DefaultConfig("https://cloud.ibm.com"), rest.NewClient())
		token, err := restClient.GetToken(tokenRequest)
		if err != nil {
			return nil, err
		}
		accessToken = token.AccessToken
	}

	configDir := config_helpers.ConfigFilePath()
	var initialConfig core_config.BXConfigData

	// check if an already existing file is present and use its information for building the config file
	if _, err := os.Stat(configDir); err == nil {
		configFile, err := os.ReadFile(configDir)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(configFile, &initialConfig)
		if err != nil {
			return nil, err
		}
	} else {
		// set the following
		//API Endpoint
		//Region
		// ... check what other config is mandatory
	}

	initialConfig.IAMToken = accessToken

	configJSON, err := json.Marshal(initialConfig)
	if err != nil {
		return nil, err
	}
	// Write access token to cache
	return configJSON, nil
}

func TryIBMCloudConfigFile() sdk.Importer {
	return importer.TryFile("~/.bluemix/config.json",
		func(ctx context.Context, contents importer.FileContents,
			in sdk.ImportInput, out *sdk.ImportAttempt) {
			var config Config
			if err := contents.ToJSON(&config); err != nil {
				out.AddError(err)
				return
			}

			if config.APIKey == "" {
				return
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.APIKey: config.APIKey,
				},
			})
		})
}

type Config struct {
	APIKey string `json:"APIKey"`
}
