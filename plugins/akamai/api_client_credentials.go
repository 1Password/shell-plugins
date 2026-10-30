package akamai

import (
	"context"
	"fmt"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
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
		DefaultProvisioner: nil,
		Importer: importer.TryAll(
			TryAkamaiConfigFile(),
		)}
}

func configFile(in sdk.ProvisionInput) ([]byte, error) {
	fields := []struct {
		name sdk.FieldName
		key  string
	}{
		{fieldname.ClientSecret, "client_secret"},
		{fieldname.Host, "host"},
		{fieldname.AccessToken, "access_token"},
		{fieldname.ClientToken, "client_token"},
	}

	var contents strings.Builder
	contents.WriteString("[default]\n")

	for _, field := range fields {
		value, ok := in.ItemFields[field.name]
		if !ok {
			continue
		}
		// INI parsers ignore surrounding whitespace, so trimming it only drops
		// harmless leading or trailing line breaks, such as from a paste.
		value = strings.TrimSpace(value)
		if strings.ContainsAny(value, "\r\n") {
			return nil, fmt.Errorf("line breaks are not allowed in the Akamai %q field", field.name)
		}
		contents.WriteString(field.key + " = " + value + "\n")
	}

	return []byte(contents.String()), nil
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
