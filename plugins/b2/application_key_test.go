package b2

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestApplicationKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, ApplicationKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.ApplicationKeyID: "00324f51a87c0d80000000001",       // 003 + 12 hex + sequence
				fieldname.ApplicationKey:   "K0038b4Q0Q+6H7ypqvt/18pGP0pEZFY", // K003 + 20 bytes base64 without padding
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"B2_ACCOUNT_ID":  "00324f51a87c0d80000000001",
					"B2_ACCOUNT_KEY": "K0038b4Q0Q+6H7ypqvt/18pGP0pEZFY",
				},
			},
		},
	})
}
