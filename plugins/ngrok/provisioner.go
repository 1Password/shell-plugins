package ngrok

import (
	"context"
	"os"
	"path/filepath"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"gopkg.in/yaml.v3"
)

const version = "2"

type fileProvisioner struct {
}

func newNgrokProvisioner() sdk.Provisioner {
	return fileProvisioner{}
}

func (f fileProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	provisionedConfigFilePath := filepath.Join(in.TempDir, "config.yml")
	var config Config
	configFilePath := processConfigFlag(out, provisionedConfigFilePath)
	if configFilePath != "" {
		existingContents, err := os.ReadFile(configFilePath)
		if err != nil {
			out.AddError(err)
			return
		}

		if err := importer.FileContents(existingContents).ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}
	}

	config.AuthToken = in.ItemFields[fieldname.Authtoken]
	config.APIKey = in.ItemFields[fieldname.APIKey]
	config.Version = version

	newContents, err := yaml.Marshal(&config)
	if err != nil {
		out.AddError(err)
		return
	}

	out.AddSecretFile(provisionedConfigFilePath, newContents)
}

func processConfigFlag(out *sdk.ProvisionOutput, newFilePath string) string {
	args := out.CommandLine
	for i, arg := range args {
		if arg == "--config" {
			if i+1 != len(args) {
				existingFilePath := args[i+1]
				args[i+1] = newFilePath
				return existingFilePath
			}
		}
	}
	args = append(args, "--config")
	args = append(args, newFilePath)
	out.CommandLine = args
	return ""
}

func (f fileProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {

}

func (f fileProvisioner) Description() string {
	return "Config file aware provisioner"
}
