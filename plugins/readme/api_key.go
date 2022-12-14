package readme

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		DocsURL:       sdk.URL("https://docs.readme.com/main/reference/intro/authentication"),
		ManagementURL: sdk.URL("https://console.readme.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Currently logged-in ReadMe user.",
				Secret:              false,
				Optional:            true,
			},
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to ReadMe.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Prefix: "rdme_",
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
			TryReadMeConfigFile(),
		)}
}

// TODO: figure out why RDME_API_KEY isn't being set as env var
var defaultEnvVarMapping = map[string]sdk.FieldName{
	"RDME_API_KEY": fieldname.APIKey,
	"RDME_EMAIL":   fieldname.Username,
	// "RDME_PROJECT": fieldname.Website,
}

func TryReadMeConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/configstore/rdme-production.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.APIKey == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			NameHint: config.Subdomain,
			Fields: map[sdk.FieldName]string{
				fieldname.Username: config.Email,
				fieldname.APIKey:   config.APIKey,
				// fieldname.Website: config.Subdomain, // TODO: figure out URL templating here
			},
		})
	})
}

type Config struct {
	APIKey    string `json:"apiKey"`
	Email     string `json:"email"`
	Subdomain string `json:"project"`
}
