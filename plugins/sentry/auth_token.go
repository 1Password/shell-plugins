package sentry

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
		DocsURL:       sdk.URL("https://docs.sentry.io/api/auth/"),
		ManagementURL: sdk.URL("https://sentry.io/settings/account/api/auth-tokens/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Sentry.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 64,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Organization,
				MarkdownDescription: "The slug of the organization to use for commands.",
				Optional:            true,
			},
			{
				Name:                fieldname.Project,
				MarkdownDescription: "The slug of the project to use for commands.",
				Optional:            true,
			},
			{
				Name:                fieldname.URL,
				MarkdownDescription: "The URL to use to connect to Sentry. This defaults to 'https://sentry.io/' but can be overridden for self-hosted instances.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TrySentryclircFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"SENTRY_AUTH_TOKEN": fieldname.Token,
	"SENTRY_ORG":        fieldname.Organization,
	"SENTRY_PROJECT":    fieldname.Project,
	"SENTRY_URL":        fieldname.URL,
}

func TrySentryclircFile() sdk.Importer {
	return importer.TryFile("~/.sentryclirc", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		fields := make(map[sdk.FieldName]string)
		for _, section := range credentialsFile.Sections() {
			if section.HasKey("token") && section.Key("token").Value() != "" {
				fields[fieldname.Token] = section.Key("token").Value()
			}

			if section.HasKey("org") && section.Key("org").Value() != "" {
				fields[fieldname.Organization] = section.Key("org").Value()
			}

			if section.HasKey("project") && section.Key("project").Value() != "" {
				fields[fieldname.Project] = section.Key("project").Value()
			}

			if section.HasKey("url") && section.Key("url").Value() != "" {
				fields[fieldname.URL] = section.Key("url").Value()
			}
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: fields,
		})
	})
}
