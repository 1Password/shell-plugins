package digitalocean

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
		DocsURL:       sdk.URL("https://docs.digitalocean.com/reference/api/create-personal-access-token/"),
		ManagementURL: sdk.URL("https://cloud.digitalocean.com/account/api/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to DigitalOcean.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 71,
					Prefix: "dop_v1_",
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		Provisioner: provision.EnvVars(map[string]string{
			fieldname.Token: "DIGITALOCEAN_ACCESS_TOKEN",
		}),
		Importer: importer.TryAllEnvVars(fieldname.Token, "DIGITALOCEAN_ACCESS_TOKEN"),
	}
}
