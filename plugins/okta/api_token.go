package okta

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIToken,
		DocsURL:       sdk.URL("https://developer.okta.com/docs/guides/create-an-api-token/main/"),
		ManagementURL: nil, // TODO: Add management URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Okta.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 42,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                FieldNameOrgURL,
				MarkdownDescription: "URL of the Okta organization to authenticate to.",
			},
		},
		Provisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryOktaConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]string{
	fieldname.Token: "OKTA_CLIENT_TOKEN",
	FieldNameOrgURL: "OKTA_CLIENT_ORGURL",
}

const FieldNameOrgURL = "Org URL"

func TryOktaConfigFile() sdk.Importer {
	// TODO: Try importing from ~/.okta/okta.yaml
	return importer.NoOp()
}
