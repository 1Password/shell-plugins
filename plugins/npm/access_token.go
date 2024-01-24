package npm

import (
	"context"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"gopkg.in/ini.v1"
)

func AccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessToken,
		DocsURL:       sdk.URL("https://docs.npmjs.com/creating-and-viewing-access-tokens"),
		ManagementURL: sdk.URL("https://www.npmjs.com/settings/<username>/tokens/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to NPM.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 36,
					Prefix: "npm_",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			TryNPMConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"NPM_CONFIG_//registry.npmjs.org/:_authToken": fieldname.Token,
}

func TryNPMConfigFile() sdk.Importer {
	return importer.TryFile("~/.npmrc", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// don't use colon as a delimiter, since it is used in the .npmrc file as a delimiter
		// between the scope, registry and configuration key
		configs, err := ini.LoadSources(ini.LoadOptions{KeyValueDelimiters: "="}, []byte(contents))
		if err != nil {
			out.AddError(err)
		}

		// sections are not supported in .npmrc
		section, err := configs.GetSection(ini.DefaultSection)
		if err != nil {
			out.AddError(err)
		}
		for _, key := range section.Keys() {
			if strings.Contains(key.Name(), "_authToken") {

				out.AddCandidate(sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: key.Value(),
					},
				})
			}
		}
	})
}
