package aws

import (
	"context"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type awsProvisioner struct {
	stsProvisioner    STSProvisioner
	envVarProvisioner provision.EnvVarProvisioner
}

func AWSProvisioner() sdk.Provisioner {
	return awsProvisioner{
		envVarProvisioner: provision.EnvVarProvisioner{
			Schema: defaultEnvVarMapping,
		},
	}
}

func (p awsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	totp, hasTotp := in.ItemFields[fieldname.OneTimePassword]
	mfaSerial, hasMFASerial := in.ItemFields[fieldname.MFASerial]

	profileWithRole := findProfileSectionWithRole(in, out)
	if len(out.Diagnostics.Errors) > 0 {
		return
	}

	if (hasTotp && hasMFASerial) || profileWithRole != nil {
		p.stsProvisioner.MFASerial = mfaSerial
		//p.stsProvisioner.MFASerial = ""
		p.stsProvisioner.TOTPCode = totp
		//p.stsProvisioner.TOTPCode = ""
		p.stsProvisioner.ProfileWithRole = profileWithRole
		p.stsProvisioner.Provision(ctx, in, out)
	} else {
		p.envVarProvisioner.Provision(ctx, in, out)
	}
}

func (p awsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	p.stsProvisioner.Deprovision(ctx, in, out)
}

func (p awsProvisioner) Description() string {
	return p.envVarProvisioner.Description()
}

func findProfileSectionWithRole(in sdk.ProvisionInput, out *sdk.ProvisionOutput) *ini.Section {
	for i, arg := range out.CommandLine {
		if arg == "--profile" {
			if i+1 == len(out.CommandLine) {
				return nil
			}
			return scanConfigFileForRole(out.CommandLine[i+1], in, out)
		}
	}

	if val := os.Getenv("AWS_PROFILE"); val != "" {
		return scanConfigFileForRole(val, in, out)
	}

	return scanConfigFileForRole("default", in, out)
}

func scanConfigFileForRole(profile string, in sdk.ProvisionInput, out *sdk.ProvisionOutput) *ini.Section {
	// Read config file from the location set in AWS_CONFIG_FILE env var or from  ~/.aws/config
	configPath := os.Getenv("AWS_CONFIG_FILE")
	if configPath != "" {
		if strings.HasPrefix(configPath, "~") {
			configPath = in.FromHomeDir(configPath[1:])
		}
	} else {
		configPath = in.FromHomeDir(".aws", "config") // default config file location
	}
	contents, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		out.AddError(fmt.Errorf("you've specified the --profile flag but no profiles to choose from were present at: %s", configPath))
		return nil
	} else if err != nil {
		out.AddError(err)
		return nil
	}

	configFile, err := importer.FileContents(contents).ToINI()
	if err != nil {
		out.AddError(err)
		return nil
	}

	if configFile != nil {
		profileSection := getConfigSectionByProfile(configFile, profile)
		if profileSection != nil {
			if profileSection.HasKey("role_arn") {
				return profileSection
			}
		}
	}

	return nil
}
