package redis

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func UserCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.UserCredentials,
		DocsURL: sdk.URL("https://redis.io/docs/ui/cli/#host-port-password-and-database"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to Redis server.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Username used to authenticate to Redis server.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "Host address for the Redis server.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Symbols:   true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Port,
				MarkdownDescription: "Port for the Redis server.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Digits: true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"REDISCLI_AUTH": fieldname.Password,
}
