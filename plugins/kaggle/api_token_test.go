package kaggle

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:    "z2pifkruzgbb17plmz2gux21fexample",
				fieldname.Username: "username",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"KAGGLE_KEY":      "z2pifkruzgbb17plmz2gux21fexample",
					"KAGGLE_USERNAME": "username",
				},
			},
		},
	})
}
