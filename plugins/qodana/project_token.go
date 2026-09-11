package qodana

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func ProjectToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIToken,
		DocsURL:       sdk.URL("https://www.jetbrains.com/help/qodana/cloud-projects.html#cloud-manage-projects"),
		ManagementURL: sdk.URL("https://qodana.cloud/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Project token used to authenticate to Qodana Cloud.",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"QODANA_TOKEN": fieldname.Token,
}
