package shopify

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCLITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "shptka_ql5v31c1kcuozk8lfdfqvuvfzexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"--password", "shptka_ql5v31c1kcuozk8lfdfqvuvfzexample"},
			},
		},
	})
}
