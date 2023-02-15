package aws

import (
	"context"
	"fmt"
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

	roleArn := findRoleArnIfSpecified(in, out)
	if len(out.Diagnostics.Errors) > 0 {
		return
	}

	if (hasTotp && hasMFASerial) || roleArn != "" {
		p.stsProvisioner.MFASerial = mfaSerial
		p.stsProvisioner.MFASerial = ""
		p.stsProvisioner.TOTPCode = totp
		p.stsProvisioner.TOTPCode = ""
		p.stsProvisioner.RoleArn = roleArn
		p.stsProvisioner.Provision(ctx, in, out)
	} else {
		p.envVarProvisioner.Provision(ctx, in, out)
	}
}

func (p awsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p awsProvisioner) Description() string {
	return p.envVarProvisioner.Description()
}

func findRoleArnIfSpecified(in sdk.ProvisionInput, out *sdk.ProvisionOutput) string {
	for i, arg := range out.CommandLine {
		if arg == "--profile" {
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
				return ""
			} else if err != nil {
				out.AddError(err)
				return ""
			}

			configFile, err := importer.FileContents(contents).ToINI()
			if err != nil {
				out.AddError(err)
				return ""
			}

			if configFile != nil {
				if i+1 == len(out.CommandLine) {
					return ""
				}
				profileSection := getConfigSectionByProfile(configFile, out.CommandLine[i+1])
				if profileSection != nil {
					if profileSection.HasKey("role_arn") {
						key, err := profileSection.GetKey("role_arn")
						if err != nil {
							out.AddError(err)
							return ""
						}

						// remove the --profile flag so the aws cli does not use it
						//out.CommandLine = out.CommandLine[0:i]
						return key.Value()
					}
				}
			}
			break
		}
	}
	return ""
}
