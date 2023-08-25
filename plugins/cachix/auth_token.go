package cachix

import (
	"context"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AuthToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AuthToken,
		DocsURL:       sdk.URL("https://docs.cachix.org/getting-started#authenticating"),
		ManagementURL: sdk.URL("https://app.cachix.org/personal-auth-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Cachix.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					// The tokens are JWT's.
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'.', '-', '_'},
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryCachixConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"CACHIX_AUTH_TOKEN": fieldname.Token,
}

func TryCachixConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/cachix/cachix.dhall", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var token string

		fileString := contents.ToString()
		// The dhall config file contains a single block with variables inside
		fileString = strings.Trim(fileString, "{}")
		// Remove spaces to make parsing simpler
		fileString = strings.Join(strings.Fields(fileString), "")
		// There should be 3 fields: authToken, hostname, binaryCaches
		keyVals := strings.SplitN(fileString, ",", 3)

		for i := range keyVals {
			kvPair := strings.Split(keyVals[i], "=")
			if len(kvPair) != 2 {
				continue
			}
			kvPair[1] = strings.Trim(kvPair[1], "\"")
			if strings.Contains(kvPair[0], "authToken") {
				token = kvPair[1]
			}
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Token: token,
			},
		})
	})
}
