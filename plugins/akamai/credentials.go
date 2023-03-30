package akamai

import (
	"context"

	"github.com/subpop/go-ini"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://techdocs.akamai.com/developer/docs/set-up-authentication-credentials"),
		ManagementURL: sdk.URL("https://control.akamai.com/apps/identity-management/#/tabs/users/list"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.ClientSecret,
				MarkdownDescription: "Client Secret used to authenticate to Akamai.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 44,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "Hostname used to authenticate to Akamai APIs.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 58,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{45, 46}, // "-" and "."
					},
				},
			},
			{
				Name:                fieldname.AccessToken,
				MarkdownDescription: "Access Token used to authenticate to Akamai.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 38,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{45}, // "-"
					},
				},
			},
			{
				Name:                fieldname.ClientToken,
				MarkdownDescription: "Client Token used to authenticate to Akamai.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 38,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{45}, // "-"
					},
				},
			},
		},
		DefaultProvisioner: provision.TempFile(configFile,
			provision.AddArgs(
				"--edgerc", "{{ .Path }}",
				"--section", "default",
			)),
		Importer: importer.TryAll(
		// TryAkamaiConfigFile(),
		)}
}

func configFile(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		Section: Section{
			ClientSecret: in.ItemFields[fieldname.ClientSecret],
			Host:         in.ItemFields[fieldname.Host],
			AccessToken:  in.ItemFields[fieldname.AccessToken],
			ClientToken:  in.ItemFields[fieldname.ClientToken],
		},
	}

	contents, err := ini.Marshal(config)
	if err != nil {
		return nil, err
	}

	return []byte(contents), nil
}

// TODO: Check if the platform stores the Credentials in a local config file, and if so,
// implement the function below to add support for importing it.
func TryAkamaiConfigFile() sdk.Importer {
	return importer.TryFile("~/.edgerc", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config.Default.ClientSecret == "" || config.Default.Host == "" || config.Default.AccessToken == "" || config.Default.ClientToken == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.ClientSecret: config.Default.ClientSecret,
		// 		fieldname.Host:         config.Default.Host,
		// 		fieldname.AccessToken:  config.Default.AccessToken,
		// 		fieldname.ClientToken:  config.Default.ClientToken,
		// 	},
		// })
	})
}

type Config struct {
	Section `ini:"default"`
}

type Section struct {
	ClientSecret string `ini:"client_secret"`
	Host         string `ini:"host"`
	AccessToken  string `ini:"access_token"`
	ClientToken  string `ini:"client_token"`
}
