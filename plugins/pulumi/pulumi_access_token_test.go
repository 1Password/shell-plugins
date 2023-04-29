package pulumi

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PulumiAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"PULUMI_ACCESS_TOKEN": "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
				},
			},
		},
		"with-backend": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
				fieldname.Host:  "https://api.pulumi.selfhosted.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"PULUMI_ACCESS_TOKEN": "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
					"PULUMI_BACKEND_URL":  "https://api.pulumi.selfhosted.com",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PulumiAccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"PULUMI_ACCESS_TOKEN": "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
					},
				},
			},
		},
		"environment-with-backend": {
			Environment: map[string]string{
				"PULUMI_ACCESS_TOKEN": "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
				"PULUMI_BACKEND_URL":  "https://api.pulumi.selfhosted.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
						fieldname.Host:  "https://api.pulumi.selfhosted.com",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.pulumi/credentials.json": plugintest.LoadFixture(t, "credentials.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
						fieldname.Host:  "http://localhost:8080",
					},
					NameHint: "localhost:8080",
				},
			},
		},
	})
}
