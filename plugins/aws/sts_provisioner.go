package aws

import (
	"context"
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type stsProvisioner struct {
	profileName string
}

func NewSTSProvisioner(profileName string) sdk.Provisioner {
	return stsProvisioner{profileName: profileName}
}

func (p stsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	awsConfig := getAWSAuthConfigurationForProfile(p.profileName, out)
	if len(out.Diagnostics.Errors) > 0 {
		return
	}

	resolveLocalAnd1PasswordConfigurations(in.ItemFields, awsConfig, out)
	if len(out.Diagnostics.Errors) > 0 {
		return
	}

	tempCredentialsProvider := p.chooseTemporaryCredentialsProvider(awsConfig, in, out)
	if len(out.Diagnostics.Errors) > 0 {
		return
	}

	tempCredentials, err := tempCredentialsProvider.Retrieve(ctx)
	if err != nil {
		out.AddError(err)
		return
	}

	out.AddEnvVar("AWS_ACCESS_KEY_ID", tempCredentials.AccessKeyID)
	out.AddEnvVar("AWS_SECRET_ACCESS_KEY", tempCredentials.SecretAccessKey)
	if tempCredentials.SessionToken != "" {
		out.AddEnvVar("AWS_SESSION_TOKEN", tempCredentials.SessionToken)
	}
	if awsConfig.Region != "" {
		out.AddEnvVar("AWS_DEFAULT_REGION", awsConfig.Region)
	}
}

func (p stsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p stsProvisioner) Description() string {
	return "Provision environment variables with temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}

// chooseTemporaryCredentialsProvider returns the aws provider that fits the scenario described by the current configuration, alongside the corresponding stsCacheWriter for encrypting temporary credentials to disk to be used in next runs.
func (p *stsProvisioner) chooseTemporaryCredentialsProvider(awsConfig *confighelpers.Config, in sdk.ProvisionInput, out *sdk.ProvisionOutput) aws.CredentialsProvider {
	unsupportedMessage := "%s is not yet supported by the AWS Shell Plugin"
	if awsConfig.HasSSOStartURL() {
		out.AddError(fmt.Errorf(unsupportedMessage, "SSO Authentication"))
	}

	if awsConfig.HasWebIdentity() {
		out.AddError(fmt.Errorf(unsupportedMessage, "Web Identity Authentication"))

	}

	if awsConfig.HasCredentialProcess() {
		out.AddError(fmt.Errorf(unsupportedMessage, "Credential Process Authentication"))

	}

	if awsConfig.HasSourceProfile() {
		out.AddError(fmt.Errorf(unsupportedMessage, "Sourcing profiles"))

	}

	if awsConfig.HasRole() {
		return NewAssumeRoleProvider(awsConfig, in, out)
	}

	if awsConfig.HasMfaSerial() && awsConfig.MfaToken != "" {
		return NewMFASessionTokenProvider(awsConfig, in, out)
	}

	return NewMasterCredentialsProvider(in.ItemFields)
}

// getAWSAuthConfigurationForProfile loads specified configurations from both config file and environment
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

	// loads configuration from both environment and config file
	configuration, err := configLoader.LoadFromProfile(profile)
	if err != nil {
		out.AddError(err)
		return nil
	}

	return configuration
}

// resolveLocalAnd1PasswordConfigurations goes over configurations present in both local settings and 1Password and resolves conflicts.
func resolveLocalAnd1PasswordConfigurations(itemFields map[sdk.FieldName]string, awsConfig *confighelpers.Config, out *sdk.ProvisionOutput) {
	mfaSerial, hasMFASerial := itemFields[fieldname.MFASerial]
	totp, hasOTP := itemFields[fieldname.OneTimePassword]
	region, hasRegion := itemFields[fieldname.DefaultRegion]

	// Give priority to the mfa serial specified in 1Password
	if hasMFASerial && awsConfig.HasMfaSerial() && awsConfig.MfaSerial != mfaSerial {
		out.AddError(fmt.Errorf("your local AWS configuration (config file or environment variable) has a different MFA serial than the one specified in 1Password"))
	}
	if !awsConfig.HasMfaSerial() {
		awsConfig.MfaSerial = mfaSerial
	}

	// Give priority to the region specified in 1Password
	if hasRegion && awsConfig.Region != "" && region != awsConfig.Region {
		out.AddError(fmt.Errorf("your local AWS configuration (config file or environment variable) has a different default region than the one specified in 1Password"))
	} else if awsConfig.Region == "" {
		awsConfig.Region = region
	}

	// only 1Password OTPs are supported
	if awsConfig.MfaToken != "" || awsConfig.MfaProcess != "" || awsConfig.MfaPromptMethod != "" {
		out.AddError(fmt.Errorf("only 1Password-backed OTP authentication is supported by the MFA worklfow of the AWS shell plugin"))
	}
	// set 1P OTP
	if hasOTP {
		awsConfig.MfaToken = totp
	}
}
