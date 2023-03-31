package akamai

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
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
			},
		},
	})
}

func TestCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, Credentials().Importer, map[string]plugintest.ImportCase{
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
