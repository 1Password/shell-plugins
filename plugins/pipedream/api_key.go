package pipedream

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
)

func APIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		DocsURL:       sdk.URL("https://pipedream.com/docs/api_key"),
		ManagementURL: sdk.URL("https://console.pipedream.com/user/security/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Pipedream.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.OrgID,
				MarkdownDescription: "OrgId for the Pipedream organization.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Length: 9,
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
		},
		DefaultProvisioner: provision.TempFile(
			pipedreamConfig,
			provision.AtFixedPath("~/.config/pipedream/config"),
		),
		Importer: importer.TryAll(
			TryPipedreamConfigFile(),
		)}
}

func TryPipedreamConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/pipedream/config", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		configFile, err := contents.ToINI()

		if err != nil {
			out.AddError(err)
			return
		}

		for _, section := range configFile.Sections() {
			fields := make(map[sdk.FieldName]string)
			if section.HasKey("api_key") && section.Key("api_key").Value() != "" {
				fields[fieldname.APIKey] = section.Key("api_key").Value()
			}
			if section.HasKey("org_id") && section.Key("org_id").Value() != "" {
				fields[fieldname.OrgID] = section.Key("org_id").Value()
			}

			if fields[fieldname.APIKey] != "" {
				out.AddCandidate(sdk.ImportCandidate{
					Fields:   fields,
					NameHint: importer.SanitizeNameHint(section.Name()),
				})
			}
		}
	})
}

type Config struct {
	APIKey string
	OrgId  string
}

func pipedreamConfig(in sdk.ProvisionInput) ([]byte, error) {
	fields := []struct {
		name sdk.FieldName
		key  string
	}{
		{fieldname.APIKey, "api_key"},
		{fieldname.OrgID, "org_id"},
	}

	var contents strings.Builder

	for _, field := range fields {
		value, ok := in.ItemFields[field.name]
		if !ok {
			continue
		}
		// INI parsers ignore surrounding whitespace, so trimming it only drops
		// harmless leading or trailing line breaks, such as from a paste.
		value = strings.TrimSpace(value)
		if strings.ContainsAny(value, "\r\n") {
			return nil, fmt.Errorf("Pipedream credential field %q cannot contain line breaks", field.name)
		}
		contents.WriteString(field.key + " = " + value + "\n")
	}

	return []byte(contents.String()), nil
}
