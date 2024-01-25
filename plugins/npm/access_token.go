package npm

import (
	"context"
	"os"
	"strings"

	"path/filepath"

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
			{
				Name:                fieldname.Organization,
				MarkdownDescription: "The organization the access token is scoped for.",
				Optional:            true,
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "The registry host for the npm packages.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.TempFile(configFile,
			provision.Filename(".npmrc"),
			provision.AddArgs(
				"--userconfig", "{{ .Path }}",
			),
		),
		Importer: importer.TryAll(
			TryNPMConfigFile(""),
			TryGlobalNPMConfigFile("NPM_CONFIG_USERCONFIG", "~"),
		),
	}
}

func configFile(in sdk.ProvisionInput) ([]byte, error) {
	contents := ""

	if org, ok := in.ItemFields[fieldname.Organization]; ok && org != "" {
		contents += "@" + strings.Trim(org, "@") + ":"
	}
	if host, ok := in.ItemFields[fieldname.Host]; ok && host != "" {
		contents += "//" + strings.Trim(host, "/") + "/:"
	}

	contents += "_authToken="

	if token, ok := in.ItemFields[fieldname.Token]; ok {
		contents += token
	}

	return []byte(contents), nil
}

func TryGlobalNPMConfigFile(env string, defaultPath string) sdk.Importer {
	path := os.Getenv(env)
	if path == "" {
		path = defaultPath
	}
	return TryNPMConfigFile(path)
}

func TryNPMConfigFile(path string) sdk.Importer {
	return importer.TryFile(filepath.Join(path, ".npmrc"), func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {

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

				keyParts := strings.Split(key.Name(), ":")

				registry := ""
				scope := ""
				hint := ""
				if len(keyParts) == 2 {
					registry = strings.Trim(keyParts[0], "/")
					hint = registry
				} else if len(keyParts) == 3 {
					registry = strings.Trim(keyParts[1], "/")
					scope = strings.Trim(keyParts[0], "@")
					hint = "@" + scope + ":" + registry
				}

				out.AddCandidate(sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:        key.Value(),
						fieldname.Host:         registry,
						fieldname.Organization: scope,
					},
					NameHint: importer.SanitizeNameHint(hint),
				})
			}
		}
	})
}
