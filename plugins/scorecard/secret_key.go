package scorecard

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func SecretKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.SecretKey,
		DocsURL:       sdk.URL("https://github.com/ossf/scorecard?tab=readme-ov-file#authentication"),
		ManagementURL: sdk.URL("https://github.com/settings/installations"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Key,
				MarkdownDescription: "RSA private key used to authenticate GitHub App for OpenSSF Scorecard. This should be a PEM-formatted private key.",
				Secret:              true,
			},
			{
				Name:                "App ID",
				MarkdownDescription: "GitHub App ID for authentication. This is a numeric ID.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Digits: true,
					},
				},
			},
			{
				Name:                "Installation ID",
				MarkdownDescription: "GitHub App Installation ID for the target repository or organization. This is a numeric ID and is required.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Digits: true,
					},
				},
			},
		},
		DefaultProvisioner: scorecardProvisioner{},
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"GITHUB_APP_KEY_PATH":         fieldname.Key,
	"GITHUB_APP_ID":               "App ID",
	"GITHUB_APP_INSTALLATION_ID": "Installation ID",
}

type scorecardProvisioner struct{}

func (p scorecardProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	// Provision the private key as a file
	if key, ok := in.ItemFields[fieldname.Key]; ok {
		keyPath := in.FromTempDir("github-app-key.pem")
		out.AddSecretFile(keyPath, []byte(key))
		out.AddEnvVar("GITHUB_APP_KEY_PATH", keyPath)
	}

	// Provision App ID and Installation ID as environment variables
	if appID, ok := in.ItemFields["App ID"]; ok && appID != "" {
		out.AddEnvVar("GITHUB_APP_ID", appID)
	}
	if installationID, ok := in.ItemFields["Installation ID"]; ok && installationID != "" {
		out.AddEnvVar("GITHUB_APP_INSTALLATION_ID", installationID)
	}
}

func (p scorecardProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Files get deleted automatically by 1Password CLI and environment variables get wiped when process exits
}

func (p scorecardProvisioner) Description() string {
	return "Provision GitHub App private key as file and IDs as environment variables for OpenSSF Scorecard"
}