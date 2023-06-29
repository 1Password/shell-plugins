package render

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Profiles map[string]Profile
}

type Profile struct {
	APIKey string `yaml:"apiKey"`
}

func APIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		DocsURL:       sdk.URL("https://render.com/docs"),
		ManagementURL: sdk.URL("https://dashboard.render.com/u/settings#api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Render.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 30,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.TempFile(renderConfig, provision.AtFixedPath("~/.render/config.yaml")),
		Importer: TryRenderConfigFile(),
	}
}

func renderConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		Profiles: map[string]Profile{
			"default": {
				APIKey: in.ItemFields[fieldname.APIKey],
			},
		},
	}
	
	contents, err := yaml.Marshal(&config)
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}

func TryRenderConfigFile() sdk.Importer {
	return importer.TryFile("~/.render/config.yaml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		for profileName, profile := range config.Profiles {
			if profile.APIKey == "" {
				continue
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.APIKey: profile.APIKey,
				},
				NameHint: importer.SanitizeNameHint(profileName),
			})
		}
	})
}

