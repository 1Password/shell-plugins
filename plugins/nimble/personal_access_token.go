package nimble

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PersonalAccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://github.com/nim-lang/nimble"),
		ManagementURL: sdk.URL("https://nimble.directory/about.html"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to GitHub.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Prefix: "ghp_",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryNimbleConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"NIMBLE_GITHUB_API_TOKEN": fieldname.Token,
}

func TryNimbleConfigFile() sdk.Importer {
	// TODO: Try importing from ~/.nimble/github_api_token file
	return importer.NoOp()
}

type Config struct {
	Token string `ini:"token"`
}
