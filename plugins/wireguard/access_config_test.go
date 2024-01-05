package wireguard

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessConfigProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessConfig().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.PrivateKey: "wvCwch6elBhsYAFlhXPj8bFEXAMPLE",
				fieldname.Address:    "10.0.0.2/32",
				fieldname.PublicKey:  "wvCwch6elBhsYAFlhXPj8bFEXAMPLE",
				fieldname.Endpoint:   "test.example.com:51820",
				fieldname.AllowedIPs: "10.0.0.0/8",
			},
			CommandLine: []string{"wg-quick"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"wg-quick", "/tmp/wg0.conf"},
				Files: map[string]sdk.OutputFile{
					"/tmp/wg0.conf": {
						Contents: []byte(plugintest.LoadFixture(t, "wg0.conf")),
					},
				},
			},
		},
	})
}

func TestAccessConfigImporter(t *testing.T) {
	expectedFields := map[sdk.FieldName]string{
		fieldname.PrivateKey: "wvCwch6elBhsYAFlhXPj8bFEXAMPLE",
		fieldname.Address:    "10.0.0.2/32",
		fieldname.PublicKey:  "wvCwch6elBhsYAFlhXPj8bFEXAMPLE",
		fieldname.Endpoint:   "test.example.com:51820",
		fieldname.AllowedIPs: "10.0.0.0/8",
	}

	plugintest.TestImporter(t, AccessConfig().Importer, map[string]plugintest.ImportCase{
		"wireguard config file": {
			Files: map[string]string{
				"/etc/wireguard/wg0.conf": plugintest.LoadFixture(t, "wg0.conf"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{Fields: expectedFields},
			},
		},
	})
}
