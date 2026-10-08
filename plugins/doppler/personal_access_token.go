package doppler

import (
	"context"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PersonalAccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://docs.doppler.com/docs/cli"),
		ManagementURL: sdk.URL("https://dashboard.doppler.com/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Personal or CLI token used to authenticate to Doppler. Personal Tokens (dp.pt.) are created in the dashboard; CLI tokens (dp.ct.) are created by `doppler login`. Both grant access to your account.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Prefix: "dp.",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'.'},
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryDopplerConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"DOPPLER_TOKEN": fieldname.Token,
}

func TryDopplerConfigFile() sdk.Importer {
	return importer.TryFile("~/.doppler/.doppler.yaml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		for _, scope := range config.Scoped {
			if !strings.HasPrefix(scope.Token, "dp.pt.") && !strings.HasPrefix(scope.Token, "dp.ct.") {
				continue
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: scope.Token,
				},
				NameHint: importer.SanitizeNameHint(scope.Project),
			})
		}
	})
}

type Config struct {
	Scoped map[string]ConfigScope `yaml:"scoped"`
}

type ConfigScope struct {
	Token   string `yaml:"token"`
	Project string `yaml:"enclave.project"`
}
