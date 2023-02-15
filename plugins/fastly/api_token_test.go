package fastly

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "4Oncq0V723ZO8HIqUgOTB77dsEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"FASTLY_API_TOKEN": "4Oncq0V723ZO8HIqUgOTB77dsEXAMPLE",
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"FASTLY_API_TOKEN": "4Oncq0V723ZO8HIqUgOTB77dsEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "4Oncq0V723ZO8HIqUgOTB77dsEXAMPLE",
					},
				},
			},
		},
		"config file on macOS": {
			OS: "darwin",
			Files: map[string]string{
				"~/Library/Application Support/fastly/config.toml": plugintest.LoadFixture(t, "config.toml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "4Oncq0V723ZO8HIqUgOTB77dsEXAMPLE",
					},
					NameHint: "first",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "NyK5NwkqXpuf74Le0omvFVUtZEXAMPLE",
					},
					NameHint: "second",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "L2IofeJGtvwy1fDSKNj5dEIRgEXAMPLE",
					},
					NameHint: "third",
				},
			},
		},
		"config file on Linux": {
			OS: "linux",
			Files: map[string]string{
				"~/.config/fastly/config.toml": plugintest.LoadFixture(t, "config.toml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "4Oncq0V723ZO8HIqUgOTB77dsEXAMPLE",
					},
					NameHint: "first",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "NyK5NwkqXpuf74Le0omvFVUtZEXAMPLE",
					},
					NameHint: "second",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "L2IofeJGtvwy1fDSKNj5dEIRgEXAMPLE",
					},
					NameHint: "third",
				},
			},
		},
	})
}
