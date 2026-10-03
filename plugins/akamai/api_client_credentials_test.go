package akamai

import (
	"fmt"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIClientCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIClientCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.ClientSecret: "abcdE23FNkBxy456z25qx9Yp5CPUxlEfQeTDkfh4QA=I",
				fieldname.Host:         "akab-lmn789n2k53w7qrs-nfkxaa4lfk3kd6ym.luna.akamaiapis.net",
				fieldname.AccessToken:  "akab-zyx987xa6osbli4k-e7jf5ikib5jknes3",
				fieldname.ClientToken:  "akab-nomoflavjuc4422e-fa2xznerxrm3teg7",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"--edgerc", "/tmp/.edgerc", "--section", "default"},
				Files: map[string]sdk.OutputFile{
					"/tmp/.edgerc": {Contents: []byte(plugintest.LoadFixture(t, ".edgerc-single"))},
				},
				Environment: map[string]string{
					"EDGERC": "/tmp/.edgerc",
				},
			},
		},
	})
}

func TestConfigFileRejectsLineBreaks(t *testing.T) {
	validFields := map[sdk.FieldName]string{
		fieldname.ClientSecret: "abcdE23FNkBxy456z25qx9Yp5CPUxlEfQeTDkfh4QA=I",
		fieldname.Host:         "akab-lmn789n2k53w7qrs-nfkxaa4lfk3kd6ym.luna.akamaiapis.net",
		fieldname.AccessToken:  "akab-zyx987xa6osbli4k-e7jf5ikib5jknes3",
		fieldname.ClientToken:  "akab-nomoflavjuc4422e-fa2xznerxrm3teg7",
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
				fields[field] = validValue + lineBreak + "debug = true"

				contents, err := configFile(sdk.ProvisionInput{ItemFields: fields})

				assert.Nil(t, contents)
				require.EqualError(t, err, fmt.Sprintf("line breaks are not allowed in the Akamai %q field", field))
				assert.NotContains(t, err.Error(), "debug = true")
			})
		}
	}
}

func TestConfigFileTrimsSurroundingWhitespace(t *testing.T) {
	validFields := map[sdk.FieldName]string{
		fieldname.ClientSecret: "abcdE23FNkBxy456z25qx9Yp5CPUxlEfQeTDkfh4QA=I",
		fieldname.Host:         "akab-lmn789n2k53w7qrs-nfkxaa4lfk3kd6ym.luna.akamaiapis.net",
		fieldname.AccessToken:  "akab-zyx987xa6osbli4k-e7jf5ikib5jknes3",
		fieldname.ClientToken:  "akab-nomoflavjuc4422e-fa2xznerxrm3teg7",
	}
	expected, err := configFile(sdk.ProvisionInput{ItemFields: validFields})
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

			contents, err := configFile(sdk.ProvisionInput{ItemFields: fields})

			require.NoError(t, err)
			assert.Equal(t, string(expected), string(contents))
		})
	}
}

func TestAPIClientCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, APIClientCredentials().Importer, map[string]plugintest.ImportCase{
		"config file with single credential": {
			Files: map[string]string{
				"~/.edgerc": plugintest.LoadFixture(t, ".edgerc-single"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "",
					Fields: map[sdk.FieldName]string{
						fieldname.ClientSecret: "abcdE23FNkBxy456z25qx9Yp5CPUxlEfQeTDkfh4QA=I",
						fieldname.Host:         "akab-lmn789n2k53w7qrs-nfkxaa4lfk3kd6ym.luna.akamaiapis.net",
						fieldname.AccessToken:  "akab-zyx987xa6osbli4k-e7jf5ikib5jknes3",
						fieldname.ClientToken:  "akab-nomoflavjuc4422e-fa2xznerxrm3teg7",
					},
				},
			},
		},
		"config file with multiple credentials": {
			Files: map[string]string{
				"~/.edgerc": plugintest.LoadFixture(t, ".edgerc-multiple"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "",
					Fields: map[sdk.FieldName]string{
						fieldname.ClientSecret: "abcdE23FNkBxy456z25qx9Yp5CPUxlEfQeTDkfh4QA=I",
						fieldname.Host:         "akab-lmn789n2k53w7qrs-nfkxaa4lfk3kd6ym.luna.akamaiapis.net",
						fieldname.AccessToken:  "akab-zyx987xa6osbli4k-e7jf5ikib5jknes3",
						fieldname.ClientToken:  "akab-nomoflavjuc4422e-fa2xznerxrm3teg7",
					},
				},
				{
					NameHint: "newcredential",
					Fields: map[sdk.FieldName]string{
						fieldname.ClientSecret: "M9XGZP/D2JedcbABC4Td8XSnHfKKIV4N5n28cj2y6zE=",
						fieldname.Host:         "akab-ip5n2k53w7nhdcxy-nflxabc432DE1ymd.luna.akamaiapis.net",
						fieldname.AccessToken:  "akab-abc77fxa6zyxi4k-e7jf5ikib5jknesc3",
						fieldname.ClientToken:  "akab-moo22awk8765efd-s2yw5zqfrx4jp57cf",
					},
				},
			},
		},
	})
}
