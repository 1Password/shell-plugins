package redis

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestSecretKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, SecretKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKey: "5v0mPzRKNcvlwRMi4CjWISt15UfCRxjcNVMPCZfDOJTZEXAMPLE",
				fieldname.SecretKey: "I2mLL1tjTKcyb5p0vWUSAcuO7XTut2QPPSSMavKQbrCEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"REDISCLOUD_ACCESS_KEY": "5v0mPzRKNcvlwRMi4CjWISt15UfCRxjcNVMPCZfDOJTZEXAMPLE",
					"REDISCLOUD_SECRET_KEY": "I2mLL1tjTKcyb5p0vWUSAcuO7XTut2QPPSSMavKQbrCEXAMPLE",
				},
			},
		},
	})
}

func TestSecretKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, SecretKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"REDISCLOUD_ACCESS_KEY": "5v0mPzRKNcvlwRMi4CjWISt15UfCRxjcNVMPCZfDOJTZEXAMPLE",
				"REDISCLOUD_SECRET_KEY": "I2mLL1tjTKcyb5p0vWUSAcuO7XTut2QPPSSMavKQbrCEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKey: "5v0mPzRKNcvlwRMi4CjWISt15UfCRxjcNVMPCZfDOJTZEXAMPLE",
						fieldname.SecretKey: "I2mLL1tjTKcyb5p0vWUSAcuO7XTut2QPPSSMavKQbrCEXAMPLE",
					},
				},
			},
		},
	})
}
