package ssh

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "ssh",
		Platform: schema.PlatformInfo{
			Name:     "SSH",
			Homepage: sdk.URL("https://www.openssh.com"),
		},
		Credentials: []schema.CredentialType{
			UserLogin(),
		},
		Executables: []schema.Executable{
			SSHCLI(),
		},
	}
}
