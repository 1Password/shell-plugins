package upcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func UserLogin() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.UserLogin,
		DocsURL:       sdk.URL("https://upcloudltd.github.io/upcloud-cli/"),
		ManagementURL: sdk.URL("https://hub.upcloud.com/people/accounts"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Username.",
				Secret:              false,
				Optional:            false,
				Composition: &schema.ValueComposition{
					Length: 20,
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password.",
				Secret:              true,
				Optional:            false,
				Composition: &schema.ValueComposition{
					Length: 36,
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
	"UPCLOUD_USERNAME": fieldname.Username,
	"UPCLOUD_PASSWORD": fieldname.Password,
}
