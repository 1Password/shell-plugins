package ngrok

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"
)

const (
	version           = "2"
	apiKeyYamlName    = "api_key"
	authTokenYamlName = "authtoken"
	versionYamlName   = "version"
	envVarAuthVersion = "v3.2.1"
)

type fileProvisioner struct {
}

func ngrokProvisioner() sdk.Provisioner {
	currentVersion, err := getNgrokVersion()
	if err != nil {
		// When ngrok version check fails for any reason,
		// use config file to provision as a fallback
		return fileProvisioner{}
	}

	// If the current ngrok CLI version is 3.2.1 or higher,
	// use environment variables to provision credentials
	//
	// semver.Compare resulting in 0 means 3.2.1 is in use
	// semver.Compare resulting in +1 means >3.2.1 is in use
	if semver.Compare(currentVersion, envVarAuthVersion) == 0 || semver.Compare(currentVersion, envVarAuthVersion) == +1 {
		newNgrokEnvVarProvisioner := ngrokEnvVarProvisioner{}
		newNgrokEnvVarProvisioner.Provision(context.Background(), sdk.ProvisionInput{}, &sdk.ProvisionOutput{})
	}

	// Otherwise use config file to provision credentials
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

		existingConfigFilePath = flagConfigFilePath
	}
	out.CommandLine = newArgs

	existingContents, err := os.ReadFile(existingConfigFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
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
		if arg == "--config" {
			if i+1 < len(args) {
				existingFilePath := args[i+1]
				args[i+1] = newFilePath
				return existingFilePath, args
			} else {
				return "", append(args, newFilePath)
			}

		} else if strings.HasPrefix(arg, "--config=") {
			existingFilePath := strings.TrimPrefix(arg, "--config=")
			args[i] = fmt.Sprintf("--config=%s", newFilePath)
			return existingFilePath, args
		}
	}
	return "", append(args, fmt.Sprintf("--config=%s", newFilePath))
}

// Get installed ngrok version and required version
func getNgrokVersion() (string, error) {
	cmd := exec.Command("ngrok", "--version")
	ngrokVersion, err := cmd.Output()
	if err != nil {
		return "", errors.New("ngrok not found")
	}

	// Example: "ngrok version 3.1.1\n" to "3.1.1\n"
	currentVersion := strings.TrimPrefix(string(ngrokVersion), "ngrok version ")

	// Example: "3.1.1\n" to "3.1.1"
	currentVersion = strings.Trim(currentVersion, "\n")

	// Example: "3.1.1" to "v3.1.1" as that's the format
	// the package semver expects
	currentVersion = "v" + currentVersion

	return currentVersion, nil
}

func (f fileProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// nothing to do here: files get deleted automatically by 1Password CLI
}

func (f fileProvisioner) Description() string {
	return "Config file aware provisioner. It will first check if an already existing config file is present."
}
