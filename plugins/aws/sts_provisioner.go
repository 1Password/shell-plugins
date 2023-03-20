package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
)

const defaultProfileName = "default"

type StsProvisioner struct {
	profileName string
}

func (p StsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	profile := p.getProfile()

	awsConfig, err := getAWSAuthConfigurationForProfile(profile, out)
	if err != nil {
		out.AddError(err)
		return
	}

	err = resolveLocalAnd1PasswordConfigurations(in.ItemFields, awsConfig)
	if err != nil {
		out.AddError(err)
		return
	}

	tempCredentialsProvider, err := p.chooseTemporaryCredentialsProvider(awsConfig, in, out)
	if err != nil {
		out.AddError(err)
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

func (p StsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p StsProvisioner) Description() string {
	return "Provision environment variables with temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}

// getProfile returns the profile to be used on this run based on specified profile information
func (p StsProvisioner) getProfile() string {
	if len(p.profileName) != 0 {
		return p.profileName
	}

	if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		return profile
	}

	return defaultProfileName
}

// chooseTemporaryCredentialsProvider returns the aws provider that fits the scenario described by the current configuration, alongside the corresponding stsCacheWriter for encrypting temporary credentials to disk to be used in next runs.
func (p *StsProvisioner) chooseTemporaryCredentialsProvider(awsConfig *confighelpers.Config, in sdk.ProvisionInput, out *sdk.ProvisionOutput) (aws.CredentialsProvider, error) {
	unsupportedMessage := "%s is not yet supported by the AWS Shell Plugin. If you would like for this feature to be supported, upvote or take on its issue: %s"
	if awsConfig.HasSSOStartURL() {
		return nil, fmt.Errorf(unsupportedMessage, "SSO Authentication", "https://github.com/1Password/shell-plugins/issues/210")
	}

	if awsConfig.HasWebIdentity() {
		return nil, fmt.Errorf(unsupportedMessage, "Web Identity Authentication", "https://github.com/1Password/shell-plugins/issues/209")

	}

	if awsConfig.HasCredentialProcess() {
		return nil, fmt.Errorf(unsupportedMessage, "Credential Process Authentication", "https://github.com/1Password/shell-plugins/issues/213")

	}

	if awsConfig.HasSourceProfile() {
		return nil, fmt.Errorf(unsupportedMessage, "Sourcing profiles", "https://github.com/1Password/shell-plugins/issues/212")

	}

	if awsConfig.HasRole() {
		return NewAssumeRoleProvider(awsConfig, in, out), nil
	}

	if awsConfig.HasMfaSerial() && awsConfig.MfaToken != "" {
		return NewMFASessionTokenProvider(awsConfig, in, out), nil
	}

	return NewAccessKeysProvider(in.ItemFields), nil
}

// getAWSAuthConfigurationForProfile loads specified configurations from both config file and environment
func getAWSAuthConfigurationForProfile(profile string, out *sdk.ProvisionOutput) (*confighelpers.Config, error) {
	// Read config file from the location set in AWS_CONFIG_FILE env var or from  ~/.aws/config
	configFile, err := confighelpers.LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	configLoader := confighelpers.ConfigLoader{
		File:          configFile,
		ActiveProfile: profile,
	}

	// loads configuration from both environment and config file
	configuration, err := configLoader.LoadFromProfile(profile)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

// resolveLocalAnd1PasswordConfigurations goes over configurations present in both local settings and 1Password and resolves conflicts using the following rules:
// - if a certain configuration is present only in 1Password, use that one.
// - if a certain configuration is present only in local configs, use that one.
// - if a certain configuration is present in both places, validate that its value is consistent between the two and use it, otherwise return an error
func resolveLocalAnd1PasswordConfigurations(itemFields map[sdk.FieldName]string, awsConfig *confighelpers.Config) error {
	mfaSerial, hasMFASerial := itemFields[fieldname.MFASerial]
	totp, hasOTP := itemFields[fieldname.OneTimePassword]
	region, hasRegion := itemFields[fieldname.DefaultRegion]

	// only 1Password OTPs are supported
	if awsConfig.MfaToken != "" || awsConfig.MfaProcess != "" || awsConfig.MfaPromptMethod != "" {
		return fmt.Errorf("only 1Password-backed OTP authentication is supported by the MFA worklfow of the AWS shell plugin")
	}
	// make sure 1Password OTP is used
	if hasOTP {
		awsConfig.MfaToken = totp
	}

	if hasMFASerial && awsConfig.HasMfaSerial() && awsConfig.MfaSerial != mfaSerial {
		return fmt.Errorf("your local AWS configuration (config file or environment variable) has a different MFA serial than the one specified in 1Password")
	} else if !awsConfig.HasMfaSerial() {
		awsConfig.MfaSerial = mfaSerial
	}

	if awsConfig.HasMfaSerial() && awsConfig.MfaToken == "" {
		return fmt.Errorf("MFA failed: an MFA serial was found but no OTP has been set up in 1Password")
	}

	if !awsConfig.HasMfaSerial() && awsConfig.MfaToken != "" {
		return fmt.Errorf("MFA failed: an OTP was found wihtout a corresponding MFA serial")
	}

	if hasRegion && awsConfig.Region != "" && region != awsConfig.Region {
		return fmt.Errorf("your local AWS configuration (config file or environment variable) has a different default region than the one specified in 1Password")
	} else if awsConfig.Region == "" {
		awsConfig.Region = region
	}

	return nil
}
