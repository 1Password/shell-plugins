package stripe

import (
	"context"
	"fmt"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/BurntSushi/toml"
)

func SecretKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.SecretKey,
		DocsURL:       sdk.URL("https://stripe.com/docs/keys"),
		ManagementURL: sdk.URL("https://dashboard.stripe.com/apikeys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Key,
				MarkdownDescription: "Key used to authenticate to Stripe.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 107,
					Prefix: "sk_",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Mode,
				MarkdownDescription: "Whether the key can be used for live or test mode.",
				Optional:            true,
			},
		},
		Provisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			importer.TryAllEnvVars(fieldname.Key, "STRIPE_SECRET_KEY"),
			TryStripeConfigFile(),
		),
	}
}

const (
	ModeLive = "live"
	ModeTest = "test"
)

var defaultEnvVarMapping = map[string]string{
	fieldname.Key: "STRIPE_API_KEY",
}

type Config struct {
	Projects map[string]ProjectConfig
}

type ProjectConfig struct {
	LiveModeAPIKey       string `toml:"live_mode_api_key"`
	LiveModeKeyExpiresAt string `toml:"live_mode_key_expires_at"`
	TestModeAPIKey       string `toml:"test_mode_api_key"`
	TestModeKeyExpiresAt string `toml:"test_mode_key_expires_at"`
}

func TryStripeConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/stripe/config.toml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		parsedFile := make(map[string]toml.Primitive)
		metaData, err := toml.Decode(string(contents), &parsedFile)
		if err != nil {
			out.AddError(err)
			return
		}

		for project, rawConfig := range parsedFile {
			var config ProjectConfig
			err = metaData.PrimitiveDecode(rawConfig, &config)
			if err != nil {
				continue // skip sections that don't define credentials
			}

			// We only support secret keys for now.
			// Support for publishable and restricted keys will be added later.
			if strings.HasPrefix(config.LiveModeAPIKey, "sk_") {
				out.AddCandidate(sdk.ImportCandidate{
					Fields: []sdk.ImportCandidateField{
						{
							Field: fieldname.Key,
							Value: config.LiveModeAPIKey,
						},
						{
							Field: fieldname.Mode,
							Value: ModeLive,
						},
					},
					NameHint: project,
				})
			}
			if strings.HasPrefix(config.TestModeAPIKey, "sk_") {
				out.AddCandidate(sdk.ImportCandidate{
					Fields: []sdk.ImportCandidateField{
						{
							Field: fieldname.Key,
							Value: config.TestModeAPIKey,
						},
						{
							Field: fieldname.Mode,
							Value: ModeTest,
						},
					},
					NameHint: fmt.Sprintf("%s – test", project),
				})
			}
		}
	})
}
