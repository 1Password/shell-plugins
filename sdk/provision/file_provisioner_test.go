package provision

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestSetOutputDirAsEnvVar(t *testing.T) {
	provisioner := TempFile(
		FieldAsFile("Config"),
		Filename("config.yaml"),
		SetOutputDirAsEnvVar("OUTPUT_DIR"),
	)

	plugintest.TestProvisioner(t, provisioner, map[string]plugintest.ProvisionCase{
		"sets_output_dir_env_var": {
			ItemFields: map[sdk.FieldName]string{"Config": ""},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{"OUTPUT_DIR": "/tmp"},
				Files: map[string]sdk.OutputFile{
					"/tmp/config.yaml": {Contents: []byte("")},
				},
			},
		},
	})
}
