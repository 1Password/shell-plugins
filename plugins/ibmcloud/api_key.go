package ibmcloud

import (
	"context"

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
		DocsURL:       sdk.URL("https://www.ibm.com/docs/en/app-connect/container?topic=servers-creating-cloud-api-key"),
		ManagementURL: sdk.URL("https://cloud.ibm.com/iam/apikeys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to IBM Cloud.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 44,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryIBMCloudConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"IBMCLOUD_API_KEY": fieldname.APIKey,
}

func TryIBMCloudConfigFile() sdk.Importer {
	return importer.TryFile("~/.bluemix/config.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.APIKey == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.APIKey: config.APIKey,
			},
		})
	})
}

type Config struct {
	APIKey string `json:"APIKey"`
}
