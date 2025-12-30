package descope

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func ManagementKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.ManagementKey,                                      
		DocsURL:       sdk.URL("https://docs.descope.com/cli/descope"),          
		ManagementURL: sdk.URL("https://app.descope.com/settings/company/managementkeys"), 
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.ProjectID,
				MarkdownDescription: "Key used to authenticate to Descope.",
				Composition: &schema.ValueComposition{
					Length: 71,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.ManagementKey,
				MarkdownDescription: "Key used to authenticate to Descope.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 71,
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
			TryDescopeConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
    "DESCOPE_PROJECT_ID":     fieldname.ProjectID,
    "DESCOPE_MANAGEMENT_KEY": fieldname.ManagementKey,
}

// TODO: Check if the platform stores the Management Key in a local config file, and if so,
// implement the function below to add support for importing it.
func TryDescopeConfigFile() sdk.Importer {
	return importer.TryFile("~/path/to/config/file.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config.Key == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.Key: config.Key,
		// 	},
		// })
	})
}

// TODO: Implement the config file schema
// type Config struct {
//	Key string
// }
