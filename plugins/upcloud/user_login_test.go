package upcloud

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, UserLogin().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Username: "theuser",
				fieldname.Password: "SuperS3cret",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"UPCLOUD_USERNAME": "theuser",
					"UPCLOUD_PASSWORD": "SuperS3cret",
				},
			},
		},
	})
}

func TestUserLoginImporter(t *testing.T) {
	plugintest.TestImporter(t, UserLogin().Importer, map[string]plugintest.ImportCase{
		"default": {
			Environment: map[string]string{
				"UPCLOUD_USERNAME": "theuser",
				"UPCLOUD_PASSWORD": "SuperS3cret",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: "theuser",
						fieldname.Password: "SuperS3cret",
					},
				},
			},
		},
	})
}
