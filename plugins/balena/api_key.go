package balena

import (
	"context"
	"fmt"

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
		DocsURL:       sdk.URL("https://www.balena.io/docs/reference/balena-cli/#login"),
		ManagementURL: sdk.URL("https://dashboard.balena-cloud.com/preferences/access-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Balena.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		Provisioner: provision.TempFile(balenaConfig),
		Importer: importer.TryAll(
			TryBalenaConfigFile(),
		)}
}
func balenaConfig(in sdk.ProvisionInput) ([]byte, error) {
	content := ""
	if APIKey, ok := in.ItemFields[fieldname.APIKey]; ok {
		content += APIKey
	}
	return []byte(content), nil
}

func TryBalenaConfigFile() sdk.Importer {
	//fmt.Println("IS IT HERE?!?")
	return importer.TryFile("~/.balena/token", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		token := string(contents)
		fields := make(map[string]string)
		fields[fieldname.APIKey] = token
		out.AddCandidate(sdk.ImportCandidate{
			Fields: fields,
		})
	})
}

func configFileEntry(key string, value string) string {
	return fmt.Sprintf("%s=%s\n", key, value)
}
