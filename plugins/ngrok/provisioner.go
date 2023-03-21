package ngrok

import (
	"context"
	"fmt"
	"github.com/1Password/shell-plugins/sdk/importer"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"gopkg.in/yaml.v3"
)

const (
	version           = "2"
	apiKeyYamlName    = "api_key"
	authTokenYamlName = "authtoken"
	versionYamlName   = "version"
)

type fileProvisioner struct {
}

func newNgrokProvisioner() sdk.Provisioner {
	return fileProvisioner{}
}

func (f fileProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	provisionedConfigFilePath := filepath.Join(in.TempDir, "config.yml")
	config := make(map[string]interface{})

	var existingConfigFilePath string
	// use default locations, depending on the OS
	switch runtime.GOOS {
	case "darwin":
		existingConfigFilePath = filepath.Join(in.HomeDir, "/Library/Application Support/ngrok/ngrok.yml")
	case "linux":
		existingConfigFilePath = filepath.Join(in.HomeDir, "/.config/ngrok/ngrok.yml")
	}

	flagConfigFilePath, newArgs := getConfigValueAndNewArgs(out.CommandLine, provisionedConfigFilePath)
	if flagConfigFilePath != "" {
		out.CommandLine = newArgs
		existingConfigFilePath = flagConfigFilePath
	}

	existingContents, err := os.ReadFile(existingConfigFilePath)
	if err != nil {
		if err != os.ErrNotExist {
			out.AddError(err)
			return
		}
	} else {
		if err := importer.FileContents(existingContents).ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}
	}

	config[authTokenYamlName] = in.ItemFields[fieldname.Authtoken]
	config[apiKeyYamlName] = in.ItemFields[fieldname.APIKey]
	config[versionYamlName] = version

	newContents, err := yaml.Marshal(&config)
	if err != nil {
		out.AddError(err)
		return
	}

	out.AddSecretFile(provisionedConfigFilePath, newContents)
}

// getConfigValueAndNewArgs returns the value of the original config flag if specified, and the arguments with the file path replaced by the newFilePath.
func getConfigValueAndNewArgs(args []string, newFilePath string) (string, []string) {
	for i, arg := range args {
		if arg == "--config" && i+1 != len(args) {
			existingFilePath := args[i+1]
			args[i+1] = newFilePath
			return existingFilePath, args

		} else if strings.HasPrefix(arg, "--config=") {
			existingFilePath := strings.TrimPrefix(arg, "--config=")
			args[i] = fmt.Sprintf("--config=%s", newFilePath)
			return existingFilePath, args
		}
	}
	return "", nil
}

func (f fileProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// nothing to do here: files get deleted automatically by 1Password CLI
}

func (f fileProvisioner) Description() string {
	return "Config file aware provisioner. It will first check if an already existing config file is present."
}
