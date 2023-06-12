package pulumi

import (
	"context"
	"net/url"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

var tokenEnvVarMapping = map[string]sdk.FieldName{
	"PULUMI_ACCESS_TOKEN": fieldname.Token,
}

func PulumiAccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://www.pulumi.com/docs/intro/pulumi-service/accounts/"),
		ManagementURL: sdk.URL("https://app.pulumi.com/account/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Pulumi.",
				Secret:              true,
				Optional:            false,
				Composition: &schema.ValueComposition{
					Length: 44,
					Prefix: "pul-",
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(tokenEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(tokenEnvVarMapping),
			TryPulumiConfigFile(),
		)}
}

const BackendOnlyCredentialName = sdk.CredentialName("Backend Endpoint")

var backendEndpointEnvVarMapping = map[string]sdk.FieldName{
	"PULUMI_BACKEND_URL": fieldname.Endpoint,
}

func PulumiBackendEndpoint() schema.CredentialType {
	return schema.CredentialType{
		Name:          BackendOnlyCredentialName,
		DocsURL:       sdk.URL("https://www.pulumi.com/docs/intro/pulumi-service/accounts/"),
		ManagementURL: sdk.URL("https://app.pulumi.com/account/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Endpoint,
				MarkdownDescription: "The URL to the Pulumi state backend. Defaults to 'https://app.pulumi.com'.",
				Secret:              false,
				Optional:            false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(backendEndpointEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(backendEndpointEnvVarMapping),
			TryPulumiConfigFile(),
		)}
}

// Duplicated from:
// https://github.com/pulumi/pulumi/blob/874a8de2dac2fe2c8138cf63f3f242660bffc738/sdk/go/common/workspace/creds.go#L112-L127
// --- START
// Account holds the information associated with a Pulumi account.
type Account struct {
	AccessToken     string    `json:"accessToken,omitempty"`     // The access token for this account.
	Username        string    `json:"username,omitempty"`        // The username for this account.
	Organizations   []string  `json:"organizations,omitempty"`   // The organizations for this account.
	LastValidatedAt time.Time `json:"lastValidatedAt,omitempty"` // The last time this token was validated.
	Insecure        bool      `json:"insecure,omitempty"`        // Allow insecure server connections when using SSL.
}

// Credentials hold the information necessary for authenticating Pulumi Cloud API requests.  It contains
// a map from the cloud API URL to the associated access token.
type Credentials struct {
	Current      string             `json:"current,omitempty"`      // the currently selected key.
	AccessTokens map[string]string  `json:"accessTokens,omitempty"` // a map of arbitrary key strings to tokens.
	Accounts     map[string]Account `json:"accounts,omitempty"`     // a map of arbitrary keys to account info.
}

// --- END

func TryPulumiConfigFile() sdk.Importer {
	return importer.TryFile("~/.pulumi/credentials.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Credentials
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		for backendUrl, accessToken := range config.AccessTokens {
			u, err := url.Parse(backendUrl)
			if err != nil {
				out.AddError(err)
				continue
			}
			if u.Scheme == "http" || u.Scheme == "https" {
				// Only add the host when it differs from the default hosted Pulumi Service
				if u.Host != "api.pulumi.com" {
					out.AddCandidate(sdk.ImportCandidate{
						Fields: map[sdk.FieldName]string{
							fieldname.Token:    accessToken,
							fieldname.Endpoint: backendUrl,
						},
						NameHint: u.Host,
					})
				} else {
					out.AddCandidate(sdk.ImportCandidate{
						Fields: map[sdk.FieldName]string{
							fieldname.Token: accessToken,
						},
					})

				}
			}

		}
	})
}
