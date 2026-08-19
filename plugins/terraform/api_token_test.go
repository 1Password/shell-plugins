package terraform

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "tlMlxpJCrwU66MOY23rPC5v8ZXe7ZjCnC5j2DaztjKdCJi20N7kTI6v86YtjOdG5t0VWYNOSnAAjvMcLsoJEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"TF_TOKEN_app_terraform_io": "tlMlxpJCrwU66MOY23rPC5v8ZXe7ZjCnC5j2DaztjKdCJi20N7kTI6v86YtjOdG5t0VWYNOSnAAjvMcLsoJEXAMPLE",
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"TF_TOKEN_app_terraform_io": "tlMlxpJCrwU66MOY23rPC5v8ZXe7ZjCnC5j2DaztjKdCJi20N7kTI6v86YtjOdG5t0VWYNOSnAAjvMcLsoJEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "tlMlxpJCrwU66MOY23rPC5v8ZXe7ZjCnC5j2DaztjKdCJi20N7kTI6v86YtjOdG5t0VWYNOSnAAjvMcLsoJEXAMPLE",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in terraform-cloud/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				// 	{
				// 		Fields: map[sdk.FieldName]string{
				// 			fieldname.Token: "tlMlxpJCrwU66MOY23rPC5v8ZXe7ZjCnC5j2DaztjKdCJi20N7kTI6v86YtjOdG5t0VWYNOSnAAjvMcLsoJEXAMPLE",
				// 		},
				// 	},
			},
		},
	})
}
