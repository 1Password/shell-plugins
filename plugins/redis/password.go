package redis

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Password() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.Password,
		DocsURL: sdk.URL("https://redis.io/docs/ui/cli/#host-port-password-and-database"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to Redis server.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "Host address for the Redis server.",
				Secret:              false,
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
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Digits: true,
					},
				},
			},
		},
		DefaultProvisioner: EnvVarFlags(flagsToProvision),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"REDISCLI_AUTH": fieldname.Password,
}

var flagsToProvision = map[string]sdk.FieldName{
	"-h": fieldname.Host,
	"-p": fieldname.Port,
}
