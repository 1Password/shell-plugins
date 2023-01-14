package argocd

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AuthToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AuthToken,
		DocsURL:       sdk.URL("https://argo-cd.readthedocs.io/en/stable/user-guide/environment-variables/"),
		ManagementURL: nil,
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AuthToken,
				MarkdownDescription: "Auth Token used to authenticate to Argo CD.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'.', '-', '_'},
					},
				},
			},
			{
				Name:                fieldname.Address,
				MarkdownDescription: "Address of the ArgoCD server without https:// prefix.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(envVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(envVarMapping),
			TryArgocdConfigFile(),
		)}
}

var envVarMapping = map[string]sdk.FieldName{
	"ARGOCD_AUTH_TOKEN": fieldname.AuthToken,
	"ARGOCD_SERVER":     fieldname.Address,
}

func TryArgocdConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/argocd/config", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		for _, context := range config.Contexts {
			fields := make(map[sdk.FieldName]string)
			nameHint := context.Name

			for _, server := range config.Servers {
				if server.ServerAddress == context.Server && server.ServerAddress != "" {
					fields[fieldname.Address] = server.ServerAddress
				}
			}

			for _, user := range config.Users {
				if user.Name == context.User && user.AuthenticationToken != "" {
					fields[fieldname.AuthToken] = user.AuthenticationToken
				}
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields:   fields,
				NameHint: importer.SanitizeNameHint(nameHint),
			})
		}
	})
}

type Config struct {
	Contexts []Context
	Servers  []Server
	Users    []User
}

type Context struct {
	Name   string `yaml:"name"`
	Server string `yaml:"server"`
	User   string `yaml:"user"`
}

type Server struct {
	ServerAddress string `yaml:"server"`
}

type User struct {
	Name                string `yaml:"name"`
	AuthenticationToken string `yaml:"auth-token"`
}
