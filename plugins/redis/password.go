package redis

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
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
				MarkdownDescription: "Password used to authenticate to Redis.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
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
