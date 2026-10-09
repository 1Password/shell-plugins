package apono

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PersonalAPIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.PersonalAPIToken,
		DocsURL: sdk.URL("https://docs.apono.io/docs/architecture-and-security/personal-api-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Personal API Token used to authenticate to Apono.",
				Secret:              true,
			},
		},
		DefaultProvisioner: aponoLoginProvisioner{},
		Importer: importer.TryAll(
			importer.MacOnly(
				TryAponoConfigFile("~/Library/Application Support/apono-cli/config.json"),
			),
			importer.LinuxOnly(
				TryAponoConfigFile("~/.config/apono-cli/config.json"),
			),
		),
	}
}

// aponoLoginProvisioner provisions the Personal API Token as the
// --personal-token flag of "apono login". The Apono CLI does not read
// credentials from environment variables or accept a token flag on other
// commands: it only takes a token at login time, after which it manages
// its own session in its config file.
type aponoLoginProvisioner struct{}

func (p aponoLoginProvisioner) Description() string {
	return "Provision Apono Personal API Token as the --personal-token flag of 'apono login'"
}

func (p aponoLoginProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	out.AddArgs("--personal-token", in.ItemFields[fieldname.Token])
}

func (p aponoLoginProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: provisioned command-line args get cleaned up automatically.
}

func TryAponoConfigFile(path string) sdk.Importer {
	return importer.TryFile(
		path,
		func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
			var config Config
			if err := contents.ToJSON(&config); err != nil {
				out.AddError(err)
				return
			}

			for profileName, profile := range config.Auth.Profiles {
				if profile.PersonalToken == "" {
					continue
				}

				out.AddCandidate(sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: profile.PersonalToken,
					},
					NameHint: importer.SanitizeNameHint(profileName),
				})
			}
		},
	)
}

type Config struct {
	Auth AuthConfig `json:"auth"`
}

type AuthConfig struct {
	ActiveProfile string                   `json:"active_profile"`
	Profiles      map[string]ProfileConfig `json:"profiles"`
}

type ProfileConfig struct {
	PersonalToken string `json:"personal_token"`
}
