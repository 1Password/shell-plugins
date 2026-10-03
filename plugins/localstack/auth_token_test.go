package localstack

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AuthToken: "ls-02b523ae-52f2-4905-b46a-0a7d7c2947aa",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"LOCALSTACK_AUTH_TOKEN": "ls-02b523ae-52f2-4905-b46a-0a7d7c2947aa",
				},
			},
		},
	})
}

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AuthToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"LOCALSTACK_AUTH_TOKEN": "ls-02b523ae-52f2-4905-b46a-0a7d7c2947aa",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AuthToken: "ls-02b523ae-52f2-4905-b46a-0a7d7c2947aa",
					},
				},
			},
		},
	})
}
