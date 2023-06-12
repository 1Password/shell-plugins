package rediscloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func RedisCloudAPIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		DocsURL:       sdk.URL("https://docs.redis.com/latest/rc/api/get-started/manage-api-keys/"),
		ManagementURL: sdk.URL("https://app.redislabs.com/#/access-management/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AccountKey,
				MarkdownDescription: "API Account key (also known as Access Key, or just API Key) to authenticate to Redis Enterprise Cloud.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.UserKey,
				MarkdownDescription: "API user key (also known as Secret Key) to authenticate to Redis Enterprise Cloud.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(envVarMappingForRedisCloud),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(envVarMappingForRedisCloud),
		)}
}

var envVarMappingForRedisCloud = map[string]sdk.FieldName{
	"REDISCLOUD_ACCESS_KEY": fieldname.AccountKey,
	"REDISCLOUD_SECRET_KEY": fieldname.UserKey,
}
