package twilio

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
		DocsURL:       sdk.URL("https://www.twilio.com/docs/glossary/what-is-an-api-key"),
		ManagementURL: sdk.URL("https://www.twilio.com/console/project/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AccountSID,
				MarkdownDescription: "Account SID used to authenticate to Twilio.",
				Composition: &schema.ValueComposition{
					Length: 34,
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
						Lowercase: true,
					},
					Prefix: "AC",
				},
			},
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Twilio.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 34,
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
						Lowercase: true,
					},
					Prefix: "SK",
				},
			},
			{
				Name:                fieldname.APISecret,
				MarkdownDescription: "API Secret used to authenticate to Twilio.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
						Lowercase: true,
					},
				},
			},
			{
				Name:                fieldname.Region,
				MarkdownDescription: "The region to use for this API Key.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryTwilioConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"TWILIO_ACCOUNT_SID": fieldname.AccountSID,
	"TWILIO_API_KEY":     fieldname.APIKey,
	"TWILIO_API_SECRET":  fieldname.APISecret,
}

func TryTwilioConfigFile() sdk.Importer {
	return importer.TryFile("~/.twilio-cli/config.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		for name, secrets := range config.Profiles {
			if secrets.AccountSid == "" || secrets.ApiKey == "" || secrets.ApiSecret == "" {
				continue
			}

			out.AddCandidate(sdk.ImportCandidate{
				NameHint: name,
				Fields: map[sdk.FieldName]string{
					fieldname.AccountSID: secrets.AccountSid,
					fieldname.APIKey:     secrets.ApiKey,
					fieldname.APISecret:  secrets.ApiSecret,
				},
			})
		}
	})
}

type Config struct {
	Profiles map[string]TwilioProfile `json:"profiles"`
}

type TwilioProfile struct {
	AccountSid string `json:"accountSid"`
	ApiKey     string `json:"apiKey"`
	ApiSecret  string `json:"apiSecret"`
}
