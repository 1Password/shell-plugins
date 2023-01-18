package digitalocean

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"DIGITALOCEAN_ACCESS_TOKEN": "dop_v1_dk98ysntlv1045mdhneztbd4o1r3q8p0tndohkpfii5m6049a8lacaq4iexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "dop_v1_dk98ysntlv1045mdhneztbd4o1r3q8p0tndohkpfii5m6049a8lacaq4iexample",
					},
				},
			},
		},
		"config file macos": {
			OS: "darwin",
			Files: map[string]string{
				"~/Library/Application Support/doctl/config.yaml": plugintest.LoadFixture(t, "config.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "dop_v1_tr33mpd5m8q9t3ncisqbceydi8dd2n60pl1yiycg97z25fkqffp8j6ycjexample",
					},
				},
			},
		},
		"config file linux": {
			OS: "linux",
			Files: map[string]string{
				digitalOceanConfigFileOnLinux(): plugintest.LoadFixture(t, "config.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "dop_v1_tr33mpd5m8q9t3ncisqbceydi8dd2n60pl1yiycg97z25fkqffp8j6ycjexample",
					},
				},
			},
		},
	})
}

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "dop_v1_dk98ysntlv1045mdhneztbd4o1r3q8p0tndohkpfii5m6049a8lacaq4iexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"DIGITALOCEAN_ACCESS_TOKEN": "dop_v1_dk98ysntlv1045mdhneztbd4o1r3q8p0tndohkpfii5m6049a8lacaq4iexample",
				},
			},
		},
	})
}
