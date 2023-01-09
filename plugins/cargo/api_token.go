package cargo

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIToken,
		DocsURL:       sdk.URL("https://doc.rust-lang.org/cargo/reference/config.html#credentials"),
		ManagementURL: sdk.URL("https://crates.io/settings/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Cargo.",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryCargoConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"CARGO_REGISTRY_TOKEN": fieldname.Token,
}

func TryCargoConfigFile() sdk.Importer {
	return importer.TryFile("~/.cargo/credentials.toml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToTOML(&config); err != nil {
			out.AddError(err)
			return
		}

		for _, configRegistry := range config.Registries {
			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: configRegistry.Token,
				},
			})

		}

	})
}

type Config struct {
	Registries []ConfigRegistry `toml:"registry"`
}

type ConfigRegistry struct {
	Token string `toml:"token"`
}
