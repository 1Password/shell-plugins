package oci

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "oci",
		Platform: schema.PlatformInfo{
			Name:     "Oracle Cloud",
			Homepage: sdk.URL("https://github.com/oracle/oci-cli"),
		},
		Credentials: []schema.CredentialType{
			AccessKey(),
		},
		Executables: []schema.Executable{
			OracleCloudCLI(),
		},
	}
}
