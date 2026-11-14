package openstack

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "openstack",
		Platform: schema.PlatformInfo{
			Name:     "OpenStack",
			Homepage: sdk.URL("https://www.openstack.org"),
		},
		Credentials: []schema.CredentialType{
			Credentials(),
		},
		Executables: []schema.Executable{
			OpenStackCLI(),
			OSC(),
		},
	}
}
