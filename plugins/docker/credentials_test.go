package docker

import (
	"os"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCredentialsProvisioner(t *testing.T) {
	/*
	credsJson, err := json.Marshal(map[string]string{"credsStore": "1password"})
	if err != nil {
		t.Fatal("Failed to marshal config file json.")
	}
	userHomeDir, err := os.UserHomeDir()
	configFileDir := filepath.Join(userHomeDir, ".docker", "config.json")
	if err != nil {
		t.Fatal("Failed to retrieve home directory.")
	}
	*/
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Username: "Stream",
				fieldname.Secret: "Stream@087",
				fieldname.Host: "https://index.docker.io/v1/",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				//TODO find a way to test for the executable file and config file
				Environment: map[string]string{
					"DOCKER_REGISTRY":"https://index.docker.io/v1/",
					"DOCKER_CREDS_USR":"Stream",
					"DOCKER_CREDS_PSW":"Stream@087",
					"PATH": os.Getenv("PATH")+":/tmp",
				},
			},
		},
	})
}

func TestCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, Credentials().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string {
				"~/.docker/config.json": plugintest.LoadFixture(t, "config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: "Stream",
						fieldname.Secret:   "Stream@087",
						fieldname.Host:     "https://index.docker.io/v1/",
					},
				},
			},
		},
	})
}
