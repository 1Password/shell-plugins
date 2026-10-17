package firectl

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			CommandLine: []string{"firectl", "users", "list"},
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey: "fw_G2bznICSywu0kUNkfqX4xpz7ECaTjrZd9EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"firectl", "users", "list", "--api-key", "fw_G2bznICSywu0kUNkfqX4xpz7ECaTjrZd9EXAMPLE"},
			},
		},
		"missing API key": {
			CommandLine: []string{"firectl", "users", "list"},
			ItemFields:  map[sdk.FieldName]string{},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"firectl", "users", "list"},
				Diagnostics: sdk.Diagnostics{
					Errors: []sdk.Error{
						{
							Message: "missing API key field",
						},
					},
				},
			},
		},
	})
}
