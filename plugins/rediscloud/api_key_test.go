package rediscloud

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestRedisCloudAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, RedisCloudAPIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccountKey: "5v0mPzRKNcvlwRMi4CjWISt15UfCRxjcNVMPCZfDOJTZEXAMPLE",
				fieldname.UserKey:    "I2mLL1tjTKcyb5p0vWUSAcuO7XTut2QPPSSMavKQbrCEXAMPLE",
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

func TestRedisCloudAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, RedisCloudAPIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"REDISCLOUD_ACCESS_KEY": "5v0mPzRKNcvlwRMi4CjWISt15UfCRxjcNVMPCZfDOJTZEXAMPLE",
				"REDISCLOUD_SECRET_KEY": "I2mLL1tjTKcyb5p0vWUSAcuO7XTut2QPPSSMavKQbrCEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccountKey: "5v0mPzRKNcvlwRMi4CjWISt15UfCRxjcNVMPCZfDOJTZEXAMPLE",
						fieldname.UserKey:    "I2mLL1tjTKcyb5p0vWUSAcuO7XTut2QPPSSMavKQbrCEXAMPLE",
					},
				},
			},
		},
	})
}
