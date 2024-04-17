package fooctl

import (
	"context"
	"fmt"

	"os"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func ConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "~/.fooctl/credentials"
	}
	return homeDir + "/.fooctl/credentials"
}

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://localhost/fooctl-cli"),
		ManagementURL: sdk.URL("https://localhost/fooctl-cli"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKeyID,
				AlternativeNames:    []string{"UUID"},
				MarkdownDescription: "UUID used to authenticate to fooctl.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Uppercase: false,
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.APIKey,
				AlternativeNames:    []string{"Token"},
				MarkdownDescription: "Token used to authenticate to fooctl.",
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
		},
		DefaultProvisioner: provision.TempFile(
			fooctlConfig,
			provision.AtFixedPath(ConfigPath()),
		),
		Importer: importer.TryAll(
			TryFooctlConfigFile(),
		)}
}

func fooctlConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		UUID:  in.ItemFields[fieldname.APIKeyID],
		Token: in.ItemFields[fieldname.APIKey],
	}

	contents := fmt.Sprintf("UUID = \"%s\"\nToken = \"%s\"", config.UUID, config.Token)
	return []byte(contents), nil
}

func TryFooctlConfigFile() sdk.Importer {
	return importer.TryFile(ConfigPath(), func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		configFile, err := contents.ToINI()

		if err != nil {
			out.AddError(err)
			return
		}

		for _, section := range configFile.Sections() {
			added := false
			fields := make(map[sdk.FieldName]string)

			addFieldIfPresent := func(key sdk.FieldName, value string) {
				if value != "" {
					fields[key] = value
					added = true
				}
			}

			addFieldIfPresent(fieldname.APIKeyID, section.Key("UUID").String())
			addFieldIfPresent(fieldname.APIKey, section.Key("Token").String())

			if added {
				out.AddCandidate(sdk.ImportCandidate{
					Fields:   fields,
					NameHint: importer.SanitizeNameHint(section.Name()),
				})
			}
		}
	})
}

type Config struct {
	UUID  string `ini:"UUID"`
	Token string `ini:"Token"`
}
