package aws

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const defaultProfileName = "default"
const defaultSessionDuration = 15 * time.Minute

type STSProvisioner struct {
	profileName        string
	newProviderFactory func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) STSProviderFactory
}

func NewSTSProvisioner(profileName string) STSProvisioner {
	return STSProvisioner{
		profileName: profileName,
		newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) STSProviderFactory {
			return &CacheProviderFactory{
				InCache:    cacheState,
				OutCache:   cacheOps,
				ItemFields: fields,
			}
		},
	}
}

// getProfile returns the profile to be used on this run based on specified profile information
func (p STSProvisioner) getProfile() (string, error) {
	if len(p.profileName) != 0 {
		return p.profileName, nil
	}

	if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		return profile, nil
	}

	return defaultProfileName, nil
}

func (p STSProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	profile, err := p.getProfile()
	if err != nil {
		out.AddError(err)
		return
	}

	awsConfig, err := getAWSAuthConfigurationForProfile(profile)
	if err != nil {
		out.AddError(err)
		return
	}

	cacheProviderFactory := p.newProviderFactory(in.Cache, out.Cache, in.ItemFields)
	tempCredentialsProvider, err := GetTemporaryCredentialsProviderForProfile(awsConfig, cacheProviderFactory, in.ItemFields)
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

func (p STSProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p STSProvisioner) Description() string {
	return "Provision environment variables with temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}

// GetTemporaryCredentialsProviderForProfile returns the aws provider that fits the scenario described by the current configuration.
func GetTemporaryCredentialsProviderForProfile(awsConfig *confighelpers.Config, providerFactory STSProviderFactory, itemFields map[sdk.FieldName]string) (aws.CredentialsProvider, error) {
	err := resolveLocalAnd1PasswordConfigurations(itemFields, awsConfig)
	if err != nil {
		return nil, err
	}

	unsupportedMessage := "%s is not yet supported by the AWS Shell Plugin. If you would like for this feature to be supported, upvote or take on its issue: %s"
	if awsConfig.HasSSOStartURL() || awsConfig.HasSSOSession() {
		return nil, fmt.Errorf(unsupportedMessage, "SSO Authentication", "https://github.com/1Password/shell-plugins/issues/210")
	}

	if awsConfig.HasWebIdentity() {
		return nil, fmt.Errorf(unsupportedMessage, "Web Identity Authentication", "https://github.com/1Password/shell-plugins/issues/209")
	}

	if awsConfig.HasCredentialProcess() {
		return nil, fmt.Errorf(unsupportedMessage, "Credential Process Authentication", "https://github.com/1Password/shell-plugins/issues/213")
	}

	var sourceCredentialsProvider aws.CredentialsProvider
	if awsConfig.HasSourceProfile() {
		sourceProfileProvider, err := GetTemporaryCredentialsProviderForProfile(awsConfig.SourceProfile, providerFactory, itemFields)
		if err != nil {
			return nil, err
		}
		sourceCredentialsProvider = sourceProfileProvider
	} else {
		sourceCredentialsProvider = providerFactory.NewAccessKeysProvider()
	}

	if awsConfig.HasRole() {
		return providerFactory.NewAssumeRoleProvider(awsConfig, sourceCredentialsProvider), nil
	}

	if awsConfig.HasMfaSerial() && awsConfig.MfaToken != "" {
		return providerFactory.NewMFASessionTokenProvider(awsConfig, sourceCredentialsProvider), nil
	}

	return sourceCredentialsProvider, nil
}

type STSProviderFactory interface {
	NewAssumeRoleProvider(awsConfig *confighelpers.Config, sourcedCredentialsProvider aws.CredentialsProvider) aws.CredentialsProvider
	NewMFASessionTokenProvider(awsConfig *confighelpers.Config, sourcedCredentialsProvider aws.CredentialsProvider) aws.CredentialsProvider
	NewAccessKeysProvider() aws.CredentialsProvider
}

// CacheProviderFactory 's functions return AWS providers that are also handling reading and writing from shell plugin's encrypted cache
type CacheProviderFactory struct {
	InCache    sdk.CacheState
	OutCache   sdk.CacheOperations
	ItemFields map[sdk.FieldName]string
}

func (m CacheProviderFactory) NewAssumeRoleProvider(awsConfig *confighelpers.Config, sourcedCredentialsProvider aws.CredentialsProvider) aws.CredentialsProvider {
	roleCacheKey := getRoleCacheKey(awsConfig.RoleARN, m.ItemFields[fieldname.AccessKeyID])
	if m.InCache.Has(roleCacheKey) {
		return NewStsCacheProvider(roleCacheKey, m.InCache)
	}

	cacheWriter := NewSTSCacheWriter(roleCacheKey, m.OutCache)

	if awsConfig.HasMfaSerial() && awsConfig.MfaToken != "" {
		return initAssumeRoleProvider(awsConfig, getSTSClient(awsConfig.Region, m.NewMFASessionTokenProvider(awsConfig, sourcedCredentialsProvider)), cacheWriter)
	}

	return initAssumeRoleProvider(awsConfig, getSTSClient(awsConfig.Region, sourcedCredentialsProvider), cacheWriter)
}

func (m CacheProviderFactory) NewMFASessionTokenProvider(awsConfig *confighelpers.Config, sourcedCredentialsProvider aws.CredentialsProvider) aws.CredentialsProvider {
	mfaCacheKey := getMfaCacheKey(m.ItemFields[fieldname.AccessKeyID])
	if m.InCache.Has(mfaCacheKey) {
		return NewStsCacheProvider(mfaCacheKey, m.InCache)
	}

	return &mfaSessionTokenProvider{
		SessionTokenProvider: confighelpers.SessionTokenProvider{
			StsClient: getSTSClient(awsConfig.Region, sourcedCredentialsProvider),
			Duration:  awsConfig.NonChainedGetSessionTokenDuration,
			Mfa:       confighelpers.NewMfa(awsConfig),
		},
		stsCacheWriter: NewSTSCacheWriter(mfaCacheKey, m.OutCache),
	}
}

func (m CacheProviderFactory) NewAccessKeysProvider() aws.CredentialsProvider {
	return accessKeysProvider{itemFields: m.ItemFields}
}

// getAWSAuthConfigurationForProfile loads specified configurations from both config file and environment
func getAWSAuthConfigurationForProfile(profile string) (*confighelpers.Config, error) {
	// Disable log output produced by AWS Vault code
	log.SetOutput(io.Discard)

	// Read config file from the location set in AWS_CONFIG_FILE env var or from  ~/.aws/config
	configFile, err := confighelpers.LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	if err = DetectSourceProfileLoop(profile, configFile); err != nil {
		return nil, err
	}

	configLoader := confighelpers.ConfigLoader{
		File:          configFile,
		ActiveProfile: profile,
		// default to minimum duration of 15 minutes for sessions. These can be overwritten by user settings
		BaseConfig: confighelpers.Config{
			AssumeRoleDuration:                defaultSessionDuration,
			NonChainedGetSessionTokenDuration: defaultSessionDuration,
		},
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
	region, hasRegularRegion := itemFields[fieldname.Region]
	defaultRegion, hasDefaultRegion := itemFields[fieldname.DefaultRegion]
	if hasDefaultRegion {
		region = defaultRegion
	}

	hasRegion := hasDefaultRegion || hasRegularRegion

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
		return fmt.Errorf(`MFA failed: MFA serial %q was detected on the associated item or in the config file for the selected profile, but no 'One-Time Password' field was found.
Learn how to add an OTP field to your item:
https://developer.1password.com/docs/cli/shell-plugins/aws/#optional-set-up-multi-factor-authentication`, awsConfig.MfaSerial)
	}

	if hasRegion && awsConfig.Region != "" && region != awsConfig.Region {
		return fmt.Errorf("your local AWS configuration (config file or environment variable) has a different default region than the one specified in 1Password")
	} else if awsConfig.Region == "" {
		awsConfig.Region = region
	}

	return nil
}

// assumeRoleProvider retrieves temporary STS credentials for an assumed role, using the plugin encrypted cache as caching layer.
type assumeRoleProvider struct {
	confighelpers.AssumeRoleProvider
	stsCacheWriter
}

func (p assumeRoleProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	credentials, err := p.AssumeRoleProvider.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	err = p.stsCacheWriter.Put(credentials)
	if err != nil {
		return aws.Credentials{}, err
	}

	return credentials, nil
}

func initAssumeRoleProvider(awsConfig *confighelpers.Config, stsClient *sts.Client, cacheWriter stsCacheWriter) *assumeRoleProvider {
	return &assumeRoleProvider{
		AssumeRoleProvider: confighelpers.AssumeRoleProvider{
			StsClient:         stsClient,
			RoleARN:           awsConfig.RoleARN,
			RoleSessionName:   awsConfig.RoleSessionName,
			ExternalID:        awsConfig.ExternalID,
			Duration:          awsConfig.AssumeRoleDuration,
			Tags:              awsConfig.SessionTags,
			TransitiveTagKeys: awsConfig.TransitiveSessionTags,
			SourceIdentity:    awsConfig.SourceIdentity,
			Mfa:               &confighelpers.Mfa{},
		},
		stsCacheWriter: cacheWriter,
	}
}

// mfaSessionTokenProvider retrieves temporary STS credentials for the MFA workflow, using the plugin encrypted cache as caching layer.
type mfaSessionTokenProvider struct {
	confighelpers.SessionTokenProvider
	stsCacheWriter
}

func (p mfaSessionTokenProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	credentials, err := p.SessionTokenProvider.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	err = p.stsCacheWriter.Put(credentials)
	if err != nil {
		return aws.Credentials{}, err
	}

	return credentials, nil
}

// stsCacheProvider retrieves temporary STS credentials from cache, given a certain key.
type stsCacheProvider struct {
	awsCacheKey string
	cache       sdk.CacheState
}

func (c stsCacheProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	var cached aws.Credentials
	if ok := c.cache.Get(c.awsCacheKey, &cached); ok {
		return cached, nil
	}

	return aws.Credentials{}, fmt.Errorf("did not find cached credentials")
}

func NewStsCacheProvider(key string, cache sdk.CacheState) aws.CredentialsProvider {
	return stsCacheProvider{
		awsCacheKey: key,
		cache:       cache,
	}
}

// stsCacheProvider retrieves the long-lived access key pair.
type accessKeysProvider struct {
	itemFields map[sdk.FieldName]string
}

func (p accessKeysProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	secret, hasSecretKey := p.itemFields[fieldname.SecretAccessKey]
	keyId, hasKeyId := p.itemFields[fieldname.AccessKeyID]

	if !hasKeyId || !hasSecretKey {
		return aws.Credentials{}, fmt.Errorf("no long lived access key pair found. Please add your Access Key Id and Secret Access Key to your 1Password item's fields")
	}

	return aws.Credentials{
		AccessKeyID:     keyId,
		SecretAccessKey: secret,
	}, nil
}

func getSTSClient(region string, credsProvider aws.CredentialsProvider) *sts.Client {
	clientConfig := aws.Config{
		Region:      region,
		Credentials: credsProvider,
	}
	return sts.NewFromConfig(clientConfig)
}

func DetectSourceProfileLoop(profile string, config *confighelpers.ConfigFile) error {
	visited := make(map[string]bool)
	sourceProfile := profile

	for sourceProfile != "" {
		if visited[sourceProfile] {
			return fmt.Errorf("infinite loop in credential configuration detected. Attempting to load from profile %q which has already been visited", sourceProfile)
		} else {
			visited[sourceProfile] = true
		}

		profileSection, ok := config.ProfileSection(sourceProfile)
		if !ok {
			return fmt.Errorf("source profile %q does not exist in your AWS config file", sourceProfile)
		}

		sourceProfile = profileSection.SourceProfile

		// profiles could source credentials from themselves. Ignore this case, as it gets gracefully handled later.
		if profileSection.Name == sourceProfile {
			break
		}
	}

	return nil
}
