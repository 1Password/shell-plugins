package stackhero

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessToken,
		DocsURL:       sdk.URL("https://www.stackhero.io/stackhero/documentations/Use-the-CLI#Non-interactive-authentication-for-scripts-CI-and-automation"),
		ManagementURL: sdk.URL("https://dashboard.stackhero.io/account/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Stackhero.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 21,
					Prefix: "usr-",
					Charset: schema.Charset{
						Lowercase: true,
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
	"STACKHERO_TOKEN": fieldname.Token,
}
