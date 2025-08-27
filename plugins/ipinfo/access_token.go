package ipinfo

import (
	"context"
	"encoding/json"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessToken,
		DocsURL:       sdk.URL("https://github.com/ipinfo/cli#installation"),
		ManagementURL: sdk.URL("https://ipinfo.io/account/token"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AccessToken,
				MarkdownDescription: "Access Token used to authenticate to IPinfo.io.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 14,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.TempFile(
			ipinfoConfig,
			provision.AtFixedPath("~/Library/Application Support/ipinfo/config.json"),
		),
		Importer: importer.TryAll(
			TryIpinfoConfigFile("~/Library/Application Support/ipinfo/config.json"),
		)}
}

func TryIpinfoConfigFile(path string) sdk.Importer {
	return importer.TryFile(
		path,
		func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
			var config Config
			if err := contents.ToJSON(&config); err != nil {
				out.AddError(err)
				return
			}

			if config.Token == "" {
				return
			}

			out.AddCandidate(
				sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessToken: config.Token,
					},
				},
			)
		},
	)
}

type Config struct {
	CacheEnabled bool   `json:"cache_enabled"`
	Token        string `json:"token"`
}

func ipinfoConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		CacheEnabled: false,
		Token:        in.ItemFields[fieldname.AccessToken],
	}
	contents, err := json.Marshal(&config)
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}
