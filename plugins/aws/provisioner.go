package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
)

type awsProvisioner struct {
	selectedProvisioner sdk.Provisioner
}

func AWSProvisioner() sdk.Provisioner {
	return awsProvisioner{}
}

func (p awsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	provisioner := p.selectCorrespondingProvisioner(in.ItemFields, out)
	if len(out.Diagnostics.Errors) > 0 {
		return
	}

	provisioner.Provision(ctx, in, out)
}

func (p awsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	p.selectedProvisioner.Deprovision(ctx, in, out)
}

func (p awsProvisioner) Description() string {
	if p.selectedProvisioner != nil {
		return p.selectedProvisioner.Description()
	}
	return "Provisions AWS with temporary credentials."
}

func (p awsProvisioner) selectCorrespondingProvisioner(itemFields map[sdk.FieldName]string, out *sdk.ProvisionOutput) sdk.Provisioner {
	awsConfig := getAWSConfiguration(out)
	if len(out.Diagnostics.Errors) > 0 {
		return nil
	}

	err := errorIfUnsupportedConfigurationIsPresent(awsConfig)
	if err != nil {
		out.AddError(err)
		return nil
	}

	selectedProvisioner := selectProvisionerBasedOnPresentValues(itemFields, awsConfig, out)

	p.selectedProvisioner = selectedProvisioner

	return selectedProvisioner
}

func getAWSConfiguration(out *sdk.ProvisionOutput) *confighelpers.Config {
	for i, arg := range out.CommandLine {
		if arg == "--profile" {
			if i+1 == len(out.CommandLine) {
				out.AddError(fmt.Errorf("--profile flag was specified without a value"))
				return nil
			}
			config := getAWSAuthConfigurationForProfile(out.CommandLine[i+1], out)
			// Remove --profile flag so aws cli doesn't attempt to assume role by itself
			out.CommandLine = removeProfileFlagFromArgs(i, out.CommandLine)
			return config
		}
	}

	if val := os.Getenv("AWS_PROFILE"); val != "" {
		return getAWSAuthConfigurationForProfile(val, out)
	}

	return getAWSAuthConfigurationForProfile("default", out)
}

func getAWSAuthConfigurationForProfile(profile string, out *sdk.ProvisionOutput) *confighelpers.Config {
	// Read config file from the location set in AWS_CONFIG_FILE env var or from  ~/.aws/config
	configFile, err := confighelpers.LoadConfigFromEnv()
	if err != nil {
		out.AddError(err)
		return nil
	}

	configLoader := confighelpers.ConfigLoader{
		File:          configFile,
		ActiveProfile: profile,
	}

	configuration, err := configLoader.LoadFromProfile(profile)
	if err != nil {
		out.AddError(err)
		return nil
	}

	return configuration
}

func selectProvisionerBasedOnPresentValues(itemFields map[sdk.FieldName]string, awsConfig *confighelpers.Config, out *sdk.ProvisionOutput) sdk.Provisioner {
	totp, hasTotp := itemFields[fieldname.OneTimePassword]
	mfaSerial, hasMFASerial := itemFields[fieldname.MFASerial]

	// Give priority to the mfa serial specified in 1Password
	if hasMFASerial && awsConfig.HasMfaSerial() && awsConfig.MfaSerial != mfaSerial {
		out.AddError(fmt.Errorf("your local AWS configuration has a different MFA serial than the one specified in 1Password"))
	}
	if !hasMFASerial {
		mfaSerial = awsConfig.MfaSerial
	}

	if awsConfig != nil && awsConfig.RoleARN != "" {
		return newAssumeRoleProvisioner(mfaSerial, totp, awsConfig)
	}

	if hasTotp && mfaSerial != "" {
		return newMFAProvisioner(mfaSerial, totp, awsConfig.Region, awsConfig.NonChainedGetSessionTokenDuration)
	}

	return provision.EnvVarProvisioner{
		Schema: defaultEnvVarMapping,
	}
}

func errorIfUnsupportedConfigurationIsPresent(awsConfig *confighelpers.Config) error {
	if awsConfig.HasSSOSession() || awsConfig.HasSSOStartURL() || awsConfig.SSORoleName != "" || awsConfig.SSOAccountID != "" || awsConfig.SSORegistrationScopes != "" || awsConfig.SSORegion != "" || awsConfig.SSOUseStdout {
		return fmt.Errorf("SSO authentication is not yet supported by the AWS shell plugin")
	}

	if awsConfig.HasWebIdentity() {
		return fmt.Errorf("web identity is not yet supported by the AWS shell plugin")
	}

	if awsConfig.HasCredentialProcess() {
		return fmt.Errorf("credential process is not yet supported by the AWS shell plugin")
	}

	if awsConfig.HasSourceProfile() {
		return fmt.Errorf("sourcing profiles is not yet supported by the AWS shell plugin")
	}

	if awsConfig.MfaToken != "" || awsConfig.MfaProcess != "" || awsConfig.MfaPromptMethod != "" {
		return fmt.Errorf("only 1Password-backed OTP authentication is supported by the MFA worklfow of the AWS shell plugin")
	}

	return nil
}

func removeProfileFlagFromArgs(argIndex int, args []string) []string {
	result := append(args[0:argIndex], args[argIndex+2:]...)
	return result
}
