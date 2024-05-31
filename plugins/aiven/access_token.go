package aiven

import (
	"context"
	"os"

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
		DocsURL:       sdk.URL("https://docs.aiven.io/docs/tools/cli"),
		ManagementURL: sdk.URL("https://console.aiven.io/profile/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AccessToken,
				MarkdownDescription: "Token used to authenticate to Aiven.io.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 392,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   false,
						Specific:  []rune{43, 47, 61}, // "+", "/" and "="
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryAivenCredentialsFileFromEnv(),
			TryAivenCredentialsFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"AIVEN_AUTH_TOKEN": fieldname.AccessToken, // https://github.com/aiven/aiven-client/blob/02765178014d5d9ded0eeec6e234606dab9f2c60/aiven/client/envdefault.py#L15
}

func TryAivenCredentialsFile() sdk.Importer {
	// aiven-credentials.json is written to ~/.config/aiven/ on Linux and OSX
	// https://github.com/aiven/aiven-client/tree/main?tab=readme-ov-file#authenticate-logins-and-tokens
	return tryFromCredentialsFile("~/.config/aiven/aiven-credentials.json")
}
func TryAivenCredentialsFileFromEnv() sdk.Importer {
	// or can be loaded from env.AIVEN_CREDENTIALS_FILE
	// we don't have to check if the path is valid here, or that the env actually returned
	// a value. The tryFromCredendialsFile will do validation of path etc
	credentialsPath := os.Getenv("AIVEN_CREDENTIALS_FILE")
	return tryFromCredentialsFile(credentialsPath)
}

func tryFromCredentialsFile(filePath string) sdk.Importer {
	return importer.TryFile(filePath, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config AivenCredentialsConfig
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.AuthToken == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.AccessToken: config.AuthToken,
			},
		})
	})
}

type AivenCredentialsConfig struct {
	AuthToken string `json:"auth_token"`
	UserEmail string `json:"user_email"`
}
