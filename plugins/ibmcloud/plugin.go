package ibmcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "ibmcloud",
		Platform: schema.PlatformInfo{
			Name:     "IBM Cloud",
			Homepage: sdk.URL("https://ibmcloud.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			IBMCloudCLI(),
		},
	}
}
