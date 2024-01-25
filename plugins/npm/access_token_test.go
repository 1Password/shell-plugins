package npm

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"/tmp/.npmrc": {Contents: []byte("//registry.npmjs.org/:_authToken=npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE")},
				},
				CommandLine: []string{"--userconfig", "/tmp/.npmrc"},
			},
		},
		"custom registry": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE",
				fieldname.Host:  "my.registry.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"/tmp/.npmrc": {Contents: []byte("//my.registry.com/:_authToken=npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE")},
				},
				CommandLine: []string{"--userconfig", "/tmp/.npmrc"},
			},
		},
		"scoped": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:        "npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE",
				fieldname.Host:         "my.registry.com",
				fieldname.Organization: "op",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"/tmp/.npmrc": {Contents: []byte("@op://my.registry.com/:_authToken=npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE")},
				},
				CommandLine: []string{"--userconfig", "/tmp/.npmrc"},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.npmrc": plugintest.LoadFixture(t, ".npmrc-default-registry"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "registry.npmjs.org",
					Fields: map[sdk.FieldName]string{
						fieldname.Token:        "npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE",
						fieldname.Host:         "registry.npmjs.org",
						fieldname.Organization: "",
					},
				},
			},
		},
	})
}
