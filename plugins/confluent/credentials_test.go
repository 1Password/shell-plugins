package confluent

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCloudCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, CloudCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Username: "test@example.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CONFLUENT_CLOUD_EMAIL": "test@example.com",
				},
			},
		},
	})
}

func TestCloudCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, CloudCredentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"CONFLUENT_CLOUD_EMAIL": "test@example.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: "test@example.com",
					},
				},
			},
		},
	})
}

//TODO: Revist once the Shell Plugins ecosystem adds support for multiple credential types per plugin
// func TestPlatformCredentialsProvisioner(t *testing.T) {
// 	plugintest.TestProvisioner(t, PlatformCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
// 		"default": {
// 			ItemFields: map[sdk.FieldName]string{
// 				fieldname.Username: "someusername",
// 			},
// 			ExpectedOutput: sdk.ProvisionOutput{
// 				Environment: map[string]string{
// 					"CONFLUENT_PLATFORM_USERNAME": "someusername",
// 				},
// 			},
// 		},
// 	})
// }

// func TestPlatformCredentialsImporter(t *testing.T) {
// 	plugintest.TestImporter(t, PlatformCredentials().Importer, map[string]plugintest.ImportCase{
// 		"environment": {
// 			Environment: map[string]string{
// 				"CONFLUENT_PLATFORM_USERNAME": "someusername",
// 			},
// 			ExpectedCandidates: []sdk.ImportCandidate{
// 				{
// 					Fields: map[sdk.FieldName]string{
// 						fieldname.Username: "someusername",
// 					},
// 				},
// 			},
// 		},
// 	})
// }
