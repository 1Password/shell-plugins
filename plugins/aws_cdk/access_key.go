package aws_cdk

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessKey,
		DocsURL:       sdk.URL("https://aws_cdk.com/docs/access_key"), // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://console.aws_cdk.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Key,
				MarkdownDescription: "Key used to authenticate to AWS CDK.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 20,
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryAWSCDKConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"AWS_CDK_KEY": fieldname.Key, // TODO: Check if this is correct
}

// TODO: Check if the platform stores the Access Key in a local config file, and if so,
// implement the function below to add support for importing it.
func TryAWSCDKConfigFile() sdk.Importer {
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
