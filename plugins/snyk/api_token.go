package snyk

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
		DocsURL:       sdk.URL("https://docs.snyk.io/snyk-api-info/authentication-for-api"),
		ManagementURL: sdk.URL("https://app.snyk.io/account"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Snyk.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		Provisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TrySnykConfigFile("$XDG_CONFIG_HOME/configstore/snyk.json"),
		)}
}

var defaultEnvVarMapping = map[string]string{
	fieldname.Token: "SNYK_TOKEN",
}

func TrySnykConfigFile(path string) sdk.Importer {
	// TODO: Import `api` field from $XDG_CONFIG_HOME/configstore/snyk.json
	return importer.NoOp()
}
