package vault

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
		DocsURL:       sdk.URL("https://developer.hashicorp.com/vault/docs/concepts/tokens"),
		ManagementURL: nil, // TODO: Add management URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to HashiCorp Vault.",
				Secret:              true,
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
			TryVaultTokenFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"VAULT_TOKEN":     fieldname.Token,
	"VAULT_ADDR":      fieldname.Address,
	"VAULT_NAMESPACE": fieldname.Namespace,
}

func TryVaultTokenFile() sdk.Importer {
	// TODO: Try importing from ~/.vault-token file
	return importer.NoOp()
}
