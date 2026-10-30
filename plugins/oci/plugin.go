package oci

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "oci",
		Platform: schema.PlatformInfo{
			Name:     "OCI",
			Homepage: sdk.URL("https://www.oracle.com/cloud/"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			OCICLI(),
		},
	}
}
