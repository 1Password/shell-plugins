package veracrypt

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestVolumePasswordProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, VolumePassword().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"password flag inserted before positional args": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "TestPassword123!",
			},
			CommandLine: []string{"-t", "--mount", "/tmp/vol", "/mnt/point"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"-t", "--mount", "-p", "TestPassword123!", "--non-interactive", "/tmp/vol", "/mnt/point"},
			},
		},
		"flags inserted before mount point": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "Secret456!",
			},
			CommandLine: []string{"--dismount", "/mnt/point"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"--dismount", "-p", "Secret456!", "--non-interactive", "/mnt/point"},
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
			CommandLine: []string{"-t", "--mount"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"-t", "--mount", "-p", "VolumePass789!", "--non-interactive"},
			},
		},
		"no command line uses flags as provided": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "MySecret123!",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"-p", "MySecret123!", "--non-interactive"},
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
