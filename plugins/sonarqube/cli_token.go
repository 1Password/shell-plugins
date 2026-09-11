package sonarqube

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func CLIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.CLIToken,
		DocsURL:       sdk.URL("https://docs.sonarsource.com/sonarqube-cli/using-sonarqube-cli/environment-variables"),
		ManagementURL: sdk.URL("https://sonarcloud.io/account/security"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "User token used to authenticate to SonarQube CLI.",
				Secret:              true,
			},
			{
				Name:                fieldname.Organization,
				MarkdownDescription: "SonarQube Cloud organization key. Use together with the token to authenticate to SonarQube Cloud.",
				Secret:              false,
				Optional:            true,
			},
			{
				Name:                fieldname.URL,
				MarkdownDescription: "Server URL. Use together with the token to authenticate to a self-hosted SonarQube instance or a specific SonarQube Cloud region.",
				Secret:              false,
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"SONARQUBE_CLI_TOKEN":  fieldname.Token,
	"SONARQUBE_CLI_ORG":    fieldname.Organization,
	"SONARQUBE_CLI_SERVER": fieldname.URL,
}
