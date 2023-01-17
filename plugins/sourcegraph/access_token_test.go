package sourcegraph

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.HostAddress: "https://sourcegraph.com",
				fieldname.Token:       "bqrv8bpqtplf7xv5lkk6oxfldtttmhzx4example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SRC_ENDPOINT":     "https://sourcegraph.com",
					"SRC_ACCESS_TOKEN": "bqrv8bpqtplf7xv5lkk6oxfldtttmhzx4example",
				},
			},
		},
	})
}
