package gitea

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"gopkg.in/yaml.v2"
)

var configPath string = "~/.config/tea/config.yml"

func PersonalAccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://gitea.com/user/settings/applications"),
		ManagementURL: nil,
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Gitea.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.HostAddress,
				MarkdownDescription: "The Gitea host to connect to, should start with https://",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Symbols:   true,
					},
					Prefix: "https://",
				},
			},
			{
				Name:                fieldname.User,
				MarkdownDescription: "Your Gitea username",
				Secret:              false,
			},
		},
		DefaultProvisioner: provision.TempFile(giteaConfig, provision.AtFixedPath(configPath)),
		Importer: importer.TryAll(
			TryGiteaConfigFile(),
		)}
}

func giteaConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		Logins: []Login{
			{
				Name:    in.ItemFields[fieldname.HostAddress],
				URL:     in.ItemFields[fieldname.HostAddress],
				Token:   in.ItemFields[fieldname.Token],
				Default: true,
				User:    in.ItemFields[fieldname.User],
			},
		},
	}
	contents, err := yaml.Marshal(&config)
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}

func TryGiteaConfigFile() sdk.Importer {
	return importer.TryFile(configPath, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		for _, login := range config.Logins {
			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.Token:       login.Token,
					fieldname.HostAddress: login.URL,
					fieldname.User:        login.User,
				},
				NameHint: importer.SanitizeNameHint(login.Name),
			})
		}
	})
}

type Login struct {
	Name              string `yaml:"name"`
	URL               string `yaml:"url"`
	Token             string `yaml:"token"`
	Default           bool   `yaml:"default"`
	SSHHost           string `yaml:"ssh_host"`
	SSHKey            string `yaml:"ssh_key"`
	Insecure          bool   `yaml:"insecure"`
	SSHCertPrincipal  string `yaml:"ssh_certificate_principal"`
	SSHAgent          bool   `yaml:"ssh_agent"`
	SSHKeyFingerprint string `yaml:"ssh_key_agent_pub"`
	SSHPassphrase     string `yaml:"-"`
	VersionCheck      bool   `yaml:"version_check"`
	User              string `yaml:"user"`
	Created           int64  `yaml:"created"`
}

type Config struct {
	Logins []Login `yaml:"logins"`
}
