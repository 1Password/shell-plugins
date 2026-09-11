package akamai

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIClientCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIClientCredentials,
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
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			TryAkamaiConfigFile(),
		)}
}

// Akamai's edgegrid library (used by the Akamai CLI and its sub-packages) checks these
// environment variables before falling back to the ~/.edgerc file. Using them avoids
// relying on command-line flag injection, which several Akamai CLI packages (e.g. cli-gtm)
// don't parse correctly when the flags are appended after the subcommand.
// See: https://github.com/akamai/AkamaiOPEN-edgegrid-golang/blob/master/pkg/edgegrid/config.go
var defaultEnvVarMapping = map[string]sdk.FieldName{
	"AKAMAI_HOST":          fieldname.Host,
	"AKAMAI_CLIENT_TOKEN":  fieldname.ClientToken,
	"AKAMAI_CLIENT_SECRET": fieldname.ClientSecret,
	"AKAMAI_ACCESS_TOKEN":  fieldname.AccessToken,
}

// Load credentials from the ~/.edgerc file.
// Import separate credentials into 1Password based on the sections in the ~/.edgerc file.
func TryAkamaiConfigFile() sdk.Importer {
	return importer.TryFile("~/.edgerc", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		configFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		for _, section := range configFile.Sections() {
			profileName := section.Name()

			fields := make(map[sdk.FieldName]string)
			if section.HasKey("client_secret") && section.Key("client_secret").Value() != "" {
				fields[fieldname.ClientSecret] = section.Key("client_secret").Value()
			}
			if section.HasKey("host") && section.Key("host").Value() != "" {
				fields[fieldname.Host] = section.Key("host").Value()
			}
			if section.HasKey("access_token") && section.Key("access_token").Value() != "" {
				fields[fieldname.AccessToken] = section.Key("access_token").Value()
			}
			if section.HasKey("client_token") && section.Key("client_token").Value() != "" {
				fields[fieldname.ClientToken] = section.Key("client_token").Value()
			}

			// add candidates that contain all required credential fields
			if fields[fieldname.ClientSecret] != "" && fields[fieldname.Host] != "" && fields[fieldname.AccessToken] != "" && fields[fieldname.ClientToken] != "" {
				out.AddCandidate(sdk.ImportCandidate{
					NameHint: importer.SanitizeNameHint(profileName),
					Fields:   fields,
				})
			}
		}
	})
}
