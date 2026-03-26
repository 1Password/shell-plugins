package pypi

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIToken,
		DocsURL:       sdk.URL("https://pypi.org/help/#apitoken"),
		ManagementURL: sdk.URL("https://pypi.org/manage/account/#api-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "API token used to authenticate to PyPI.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Prefix: "pypi-",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'-', '_'},
					},
				},
			},
		},
		DefaultProvisioner: PyPIToolProvisioner("TWINE_USERNAME", "TWINE_PASSWORD"),
		Importer: importer.TryAll(
			importer.TryAllEnvVars(fieldname.Token, "TWINE_PASSWORD", "FLIT_PASSWORD", "HATCH_INDEX_AUTH"),
			TryPyPIRCFile(),
		),
	}
}
