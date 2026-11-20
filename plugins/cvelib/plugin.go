package cvelib

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "cvelib",
		Platform: schema.PlatformInfo{
			Name:     "CVE Services",
			Homepage: sdk.URL("https://www.cve.org/AllResources/CveServices"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			CVEServicesAPICLI(),
		},
	}
}
