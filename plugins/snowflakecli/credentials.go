package snowflakecli

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

var (
	fieldPrivateKeyPassphrase = sdk.FieldName("Private Key Passphrase")
	fieldRole                 = sdk.FieldName("Role")
	fieldWarehouse            = sdk.FieldName("Warehouse")
	fieldSchema               = sdk.FieldName("Schema")
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://docs.snowflake.com/en/developer-guide/snowflake-cli/connecting/configure-connections#use-a-private-key-file-for-authentication"),
		ManagementURL: sdk.URL("https://docs.snowflake.com/en/user-guide/key-pair-auth"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Account,
				MarkdownDescription: "Snowflake account identifier, e.g. `myorg-myaccount`.",
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Snowflake username.",
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.PrivateKey,
				MarkdownDescription: "Contents of the private key in PKCS#8 PEM format (e.g. `rsa_key.p8`), from `-----BEGIN PRIVATE KEY-----` to `-----END PRIVATE KEY-----`. To import it from a file, run `op plugin init snow` with `SNOWFLAKE_PRIVATE_KEY_FILE` set to the path of the key.",
				Secret:              true,
			},
			{
				Name:                fieldPrivateKeyPassphrase,
				MarkdownDescription: "Passphrase of the private key, if it is encrypted.",
				Secret:              true,
				Optional:            true,
			},
			{
				Name:                fieldRole,
				MarkdownDescription: "Role to use for the connection.",
				Optional:            true,
			},
			{
				Name:                fieldWarehouse,
				MarkdownDescription: "Warehouse to use for the connection.",
				Optional:            true,
			},
			{
				Name:                fieldname.Database,
				MarkdownDescription: "Database to use for the connection.",
				Optional:            true,
			},
			{
				Name:                fieldSchema,
				MarkdownDescription: "Schema to use for the connection.",
				Optional:            true,
			},
		},
		DefaultProvisioner: snowProvisioner(),
		Importer: importer.TryAll(
			TrySnowflakeEnvVars(),
			TrySnowflakeConfigFiles(),
		),
	}
}
