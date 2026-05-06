package veracrypt

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestVolumePasswordProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, VolumePassword().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"password flag injection": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "TestPassword123!",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"-p", "TestPassword123!", "--non-interactive"},
			},
		},
		"includes non-interactive flag": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "Secret456!",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"-p", "Secret456!", "--non-interactive"},
			},
		},
		"empty password returns error": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{},
				Diagnostics: sdk.Diagnostics{
					Errors: []sdk.Error{
						{Message: "password is required"},
					},
				},
			},
		},
		"volume field stored but not in command line": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "VolumePass789!",
				"Volume":           "/path/to/volume.tc",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"-p", "VolumePass789!", "--non-interactive"},
			},
		},
	})
}

func TestVolumePasswordImporter(t *testing.T) {
	plugintest.TestImporter(t, VolumePassword().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"VERACRYPT_PASSWORD": "TestPassword123!",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Password: "TestPassword123!",
					},
				},
			},
		},
	})
}
