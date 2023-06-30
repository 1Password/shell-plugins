package appwrite

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"temp file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey:   "zsaugacpwq6k54nnbdbmh1cys98u2a32qqkacma2ioxn1e2j6eyrk9urom0vzcvm6qbbm8s6l4xbm86n37foauiqba9tlcvohuoz87j7nwvpob5wr71k58i105fn39a10vj7ob84opwf1vrfat3m8konch7xxy2z2dh1ykohdbef7xgmvtn82lebe4mzmfzoylqy4jslrok11zbjtmd6xs84ukd7b1k9ofyuanvinmlhkgua32p5x0gqbexample",
				fieldname.Endpoint: "http://localhost/v1",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					ConfigPath(): {
						Contents: []byte(plugintest.LoadFixture(t, "import_prefs.json")),
					},
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"Appwrite prefs file": {
			Files: map[string]string{
				ConfigPath(): plugintest.LoadFixture(t, "import_prefs.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey:   "zsaugacpwq6k54nnbdbmh1cys98u2a32qqkacma2ioxn1e2j6eyrk9urom0vzcvm6qbbm8s6l4xbm86n37foauiqba9tlcvohuoz87j7nwvpob5wr71k58i105fn39a10vj7ob84opwf1vrfat3m8konch7xxy2z2dh1ykohdbef7xgmvtn82lebe4mzmfzoylqy4jslrok11zbjtmd6xs84ukd7b1k9ofyuanvinmlhkgua32p5x0gqbexample",
						fieldname.Endpoint: "http://localhost/v1",
					},
				},
			},
		},
	})
}
