package qdrantcloud

import (
	"context"
	"os"

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
		DocsURL:       sdk.URL("https://qdrant.tech/documentation/cloud-api/"),
		ManagementURL: sdk.URL("https://cloud.qdrant.io"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.ManagementKey,
				MarkdownDescription: "Cloud Management Key used to authenticate to the Qdrant Cloud API.",
				Secret:              true,
			},
			{
				Name:                fieldname.AccountID,
				MarkdownDescription: "ID of the Qdrant Cloud account to manage.",
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryQdrantCloudConfigFileFromEnv(),
			TryQdrantCloudConfigFile("~/.config/qcloud/config.yaml"),
			TryQdrantCloudConfigFile("~/.config/qcloud/config.yml"),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"QDRANT_CLOUD_API_KEY":    fieldname.ManagementKey,
	"QDRANT_CLOUD_ACCOUNT_ID": fieldname.AccountID,
}

// TryQdrantCloudConfigFileFromEnv imports from the config file that QDRANT_CLOUD_CONFIG points qcloud to, if set.
func TryQdrantCloudConfigFileFromEnv() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		if path := os.Getenv("QDRANT_CLOUD_CONFIG"); path != "" {
			TryQdrantCloudConfigFile(path)(ctx, in, out)
		}
	}
}

// TryQdrantCloudConfigFile imports every context that `qcloud context set` stored with a plaintext api_key.
// Contexts using api_key_command are skipped, since their key doesn't live in the config file.
func TryQdrantCloudConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		for _, c := range config.Contexts {
			if c.APIKey == "" {
				continue
			}

			fields := map[sdk.FieldName]string{
				fieldname.ManagementKey: c.APIKey,
			}
			if c.AccountID != "" {
				fields[fieldname.AccountID] = c.AccountID
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields:   fields,
				NameHint: importer.SanitizeNameHint(c.Name),
			})
		}
	})
}

type Config struct {
	Contexts []ConfigContext `yaml:"contexts"`
}

type ConfigContext struct {
	Name      string `yaml:"name"`
	APIKey    string `yaml:"api_key"`
	AccountID string `yaml:"account_id"`
}
