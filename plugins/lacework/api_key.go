package lacework

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/BurntSushi/toml"
)

func APIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		DocsURL:       sdk.URL("https://docs.lacework.com/console/api-access-keys"),
		ManagementURL: sdk.URL("https://login.lacework.net/ui?redirectUrl=/investigation/settings/apikeys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Account,
				MarkdownDescription: "The subdomain used to access your Lacework account.",
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Specific:  []rune{'.'},
					},
				},
			},
			{
				Name:                fieldname.APIKeyID,
				MarkdownDescription: "The API key used to authenticate to Lacework.",
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
						Specific:  []rune{'-'},
					},
				},
			},
			{
				Name:                fieldname.APISecret,
				MarkdownDescription: "The API secret used to authenticate to Lacework.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 33,
					Prefix: "_",
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
			TryLaceworkConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"LW_ACCOUNT":    fieldname.Account,
	"LW_API_KEY":    fieldname.APIKeyID,
	"LW_API_SECRET": fieldname.APISecret,
}

func TryLaceworkConfigFile() sdk.Importer {
	return importer.TryFile("~/.lacework.toml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		parsedFile := make(map[string]toml.Primitive)
		metaData, err := toml.Decode(string(contents), &parsedFile)
		if err != nil {
			out.AddError(err)
			return
		}

		for profileName, rawConfig := range parsedFile {
			var config ProfileConfig
			err = metaData.PrimitiveDecode(rawConfig, &config)
			if err != nil {
				continue // skip sections that don't define credentials
			}

			if profileName == "default" {
				profileName = ""
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.Account:   config.Account,
					fieldname.APIKeyID:  config.APIKey,
					fieldname.APISecret: config.APISecret,
				},
				NameHint: profileName,
			})
		}
	})
}

type Config struct {
	Profiles map[string]ProfileConfig
}

type ProfileConfig struct {
	Account   string `toml:"account"`
	APIKey    string `toml:"api_key"`
	APISecret string `toml:"api_secret"`
}
