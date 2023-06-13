package redis

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestUserCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, redisProvisioner(), map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "pjtxpc2gaddifapjvalggspojexample",
				fieldname.Username: "example",
				fieldname.Host:     "127.0.0.1",
				fieldname.Port:     "6379",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"REDISCLI_AUTH": "pjtxpc2gaddifapjvalggspojexample",
				},
				CommandLine: []string{"-p", "6379", "-h", "127.0.0.1", "--user", "example"},
			},
		},
	})
}
func TestDefaultUserCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, UserCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "pjtxpc2gaddifapjvalggspojexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"REDISCLI_AUTH": "pjtxpc2gaddifapjvalggspojexample",
				},
			},
		},
	})
}

func TestUserCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, UserCredentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"REDISCLI_AUTH": "pjtxpc2gaddifapjvalggspojexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Password: "pjtxpc2gaddifapjvalggspojexample",
					},
				},
			},
		},
	})
}
