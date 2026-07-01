package openbao

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AuthToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AuthToken,
		DocsURL:       sdk.URL("https://openbao.org/docs/concepts/tokens/"),
		ManagementURL: nil,
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to OpenBao.",
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
				Name:                fieldname.Address,
				MarkdownDescription: "Default address of the Vault server to use for this auth token.",
				Optional:            true,
			},
			{
				Name:                fieldname.Namespace,
				MarkdownDescription: "Default namespace to use for this auth token.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryOpenBaoConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"BAO_TOKEN":     fieldname.Token,
	"BAO_ADDR":      fieldname.Address,
	"BAO_NAMESPACE": fieldname.Namespace,
}

func TryOpenBaoConfigFile() sdk.Importer {
	return importer.NoOp()
}

