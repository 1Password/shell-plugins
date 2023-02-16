package ibmcloud

import (
	"context"
	"encoding/json"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/IBM-Cloud/bluemix-go/endpoints"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/authentication/iam"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/config_helpers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	"os"
	"path/filepath"
)

const cloudEndpoint = "https://cloud.ibm.com"

type ibmProvisioner struct {
	provision.FileProvisioner
}

func IBMProvisioner() sdk.Provisioner {
	return ibmProvisioner{}
}

func (p ibmProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	sess, err := session.New()
	if err != nil {
		out.AddError(err)
		return
	}
	config := sess.Config

	var outputConfig core_config.BXConfigData

	configPath := filepath.Join(config_helpers.ConfigDir(), "config.json")
	if _, err := os.Stat(configPath); err == nil {
		configFile, err := os.ReadFile(configPath)
		if err != nil {
			out.AddError(err)
			return
		}

		err = json.Unmarshal(configFile, &outputConfig)
		if err != nil {
			out.AddError(err)
			return
		}
	} else {
		endpointProvider := endpoints.NewEndpointLocator(config.Region, config.Visibility, config.EndpointsFile)
		IAMEndpoint, err := endpointProvider.IAMEndpoint()
		if err != nil {
			out.AddError(err)
			return
		}
		outputConfig.IAMEndpoint = IAMEndpoint
		outputConfig.APIEndpoint = cloudEndpoint
		outputConfig.Region = config.Region
	}

	if ok := in.Cache.Get("ibm-access-token", config.IAMAccessToken); !ok {
		tokenRequest := iam.APIKeyTokenRequest(in.ItemFields[fieldname.APIKey])
		restClient := iam.NewClient(iam.DefaultConfig(outputConfig.IAMEndpoint), rest.NewClient())
		token, err := restClient.GetToken(tokenRequest)
		if err != nil {
			out.AddError(err)
			return
		}
		config.IAMAccessToken = token.AccessToken

		err = out.Cache.Put("ibm-access-token", token.AccessToken, token.Expiry)
		if err != nil {
			out.AddError(err)
			return
		}
	}
	outputConfig.IAMToken = config.IAMAccessToken

	jsonConfig, err := json.Marshal(outputConfig)
	if err != nil {
		out.AddError(err)
		return
	}

	path := filepath.Join(in.TempDir, "config.json")
	out.AddEnvVar("IBMCLOUD_CONFIG_HOME", in.TempDir)

	out.AddSecretFile(path, jsonConfig)
}

func (p ibmProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: files get wiped automatically when the process exits.
}

func (p ibmProvisioner) Description() string {
	return "File provisioner for the IBM Cloud."
}
