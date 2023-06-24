package vercel

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIToken,
		DocsURL:       sdk.URL("https://vercel.com/docs/rest-api#authentication"),
		ManagementURL: sdk.URL("https://vercel.com/account/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Vercel.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 24,
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: vercelProvisioner{},
		Importer: importer.TryAll(
			importer.MacOnly(TryVercelConfigFile("~/Library/Application Support/com.vercel.cli/auth.json")),
			importer.LinuxOnly(TryVercelConfigFile("~/.config/com.vercel.cli/auth.json")),
		),
	}
}

type vercelProvisioner struct{}

func (v vercelProvisioner) Description() string {
	return "Vercel cli token provisioner"
}

func (v vercelProvisioner) Provision(ctx context.Context, input sdk.ProvisionInput, output *sdk.ProvisionOutput) {
	output.AddArgs("--token", input.ItemFields[fieldname.Token])
}

func (v vercelProvisioner) Deprovision(ctx context.Context, input sdk.DeprovisionInput, output *sdk.DeprovisionOutput) {
	// No-op
}

func TryVercelConfigFile(path string) sdk.Importer {
	return importer.TryFile(
		path,
		func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
			var config Config
			if err := contents.ToJSON(&config); err != nil {
				out.AddError(err)
				return
			}

			if config.Token == "" {
				return
			}

			out.AddCandidate(
				sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: config.Token,
					},
				},
			)
		},
	)
}

type Config struct {
	Token string `json:"token"`
}
