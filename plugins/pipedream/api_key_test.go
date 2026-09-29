package pipedream

import (
	"fmt"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"config file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey: "ugvfxesz62ycsl42z49c0t1hjexample",
				fieldname.OrgID:  "YbEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"~/.config/pipedream/config": {
						Contents: []byte(plugintest.LoadFixture(t, "provision")),
					},
				},
			},
		},
	})
}

func TestPipedreamConfigRejectsLineBreaks(t *testing.T) {
	validFields := map[sdk.FieldName]string{
		fieldname.APIKey: "ugvfxesz62ycsl42z49c0t1hjexample",
		fieldname.OrgID:  "YbEXAMPLE",
	}

	for field, validValue := range validFields {
		for name, lineBreak := range map[string]string{
			"line feed":       "\n",
			"carriage return": "\r",
		} {
			t.Run(fmt.Sprintf("%s/%s", field, name), func(t *testing.T) {
				fields := make(map[sdk.FieldName]string, len(validFields))
				for key, value := range validFields {
					fields[key] = value
				}
				fields[field] = validValue + lineBreak + "[other]"

				contents, err := pipedreamConfig(sdk.ProvisionInput{ItemFields: fields})

				assert.Nil(t, contents)
				require.EqualError(t, err, fmt.Sprintf("line breaks are not allowed in the Pipedream %q field", field))
				assert.NotContains(t, err.Error(), "[other]")
			})
		}
	}
}

func TestPipedreamConfigTrimsSurroundingWhitespace(t *testing.T) {
	validFields := map[sdk.FieldName]string{
		fieldname.APIKey: "ugvfxesz62ycsl42z49c0t1hjexample",
		fieldname.OrgID:  "YbEXAMPLE",
	}
	expected, err := pipedreamConfig(sdk.ProvisionInput{ItemFields: validFields})
	require.NoError(t, err)

	for name, pad := range map[string]func(string) string{
		"trailing line feed":       func(v string) string { return v + "\n" },
		"trailing CRLF":            func(v string) string { return v + "\r\n" },
		"leading line feed":        func(v string) string { return "\n" + v },
		"surrounding spaces, tabs": func(v string) string { return " \t" + v + "\t " },
	} {
		t.Run(name, func(t *testing.T) {
			fields := make(map[sdk.FieldName]string, len(validFields))
			for key, value := range validFields {
				fields[key] = pad(value)
			}

			contents, err := pipedreamConfig(sdk.ProvisionInput{ItemFields: fields})

			require.NoError(t, err)
			assert.Equal(t, string(expected), string(contents))
		})
	}
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.config/pipedream/config": plugintest.LoadFixture(t, "import"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "ugvfxesz62ycsl42z49c0t1hjexample",
						fieldname.OrgID:  "YbEXAMPLE",
					},
					NameHint: "DEFAULT",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "5puf32rvhkz83c6oj4wpxvaniexample",
						fieldname.OrgID:  "KVEXAMPLE",
					},
					NameHint: "first",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "lgx1amb0qf7mjy6y7nkgfc3x9example",
					},
					NameHint: "second",
				},
			},
		},
	})
}
