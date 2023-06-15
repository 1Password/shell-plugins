package aws

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/skratchdot/open-golang/open"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const defaultProfileName = "default"

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

	err = resolveLocalAnd1PasswordConfigurations(in.ItemFields, awsConfig)
	if err != nil {
		out.AddError(err)
		return
	}

	cacheProviderFactory := p.newProviderFactory(in.Cache, out.Cache, in.ItemFields)
	tempCredentialsProvider, err := ChooseTemporaryCredentialsProvider(awsConfig, cacheProviderFactory)
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

// ChooseTemporaryCredentialsProvider returns the aws provider that fits the scenario described by the current configuration.
func ChooseTemporaryCredentialsProvider(awsConfig *confighelpers.Config, providerFactory STSProviderFactory) (aws.CredentialsProvider, error) {
	unsupportedMessage := "%s is not yet supported by the AWS Shell Plugin. If you would like for this feature to be supported, upvote or take on its issue: %s"
	if awsConfig.HasSSOStartURL() || awsConfig.HasSSOSession() {
		return providerFactory.NewSSOProvider(awsConfig), nil
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
		return providerFactory.NewAssumeRoleProvider(awsConfig), nil
	}

	if awsConfig.HasMfaSerial() && awsConfig.MfaToken != "" {
		return providerFactory.NewMFASessionTokenProvider(awsConfig), nil
	}

	return providerFactory.NewAccessKeysProvider(), nil
}

type STSProviderFactory interface {
	NewAssumeRoleProvider(awsConfig *confighelpers.Config) aws.CredentialsProvider
	NewMFASessionTokenProvider(awsConfig *confighelpers.Config) aws.CredentialsProvider
	NewAccessKeysProvider() aws.CredentialsProvider
	NewSSOProvider(awsConfig *confighelpers.Config) aws.CredentialsProvider
}

// CacheProviderFactory 's functions return AWS providers that are also handling reading and writing from shell plugin's encrypted cache
type CacheProviderFactory struct {
	InCache    sdk.CacheState
	OutCache   sdk.CacheOperations
	ItemFields map[sdk.FieldName]string
}

func (m CacheProviderFactory) NewAssumeRoleProvider(awsConfig *confighelpers.Config) aws.CredentialsProvider {
	roleCacheKey := getRoleCacheKey(awsConfig.RoleARN, m.ItemFields[fieldname.AccessKeyID])
	if m.InCache.Has(roleCacheKey) {
		return NewStsCacheProvider(roleCacheKey, m.InCache)
	}

	cacheWriter := NewSTSCacheWriter(roleCacheKey, m.OutCache)

	if awsConfig.HasMfaSerial() && awsConfig.MfaToken != "" {
		return initAssumeRoleProvider(awsConfig, getSTSClient(awsConfig.Region, m.NewMFASessionTokenProvider(awsConfig)), cacheWriter)
	}

	return initAssumeRoleProvider(awsConfig, getSTSClient(awsConfig.Region, m.NewAccessKeysProvider()), cacheWriter)
}

func (m CacheProviderFactory) NewMFASessionTokenProvider(awsConfig *confighelpers.Config) aws.CredentialsProvider {
	mfaCacheKey := getMfaCacheKey(m.ItemFields[fieldname.AccessKeyID])
	if m.InCache.Has(mfaCacheKey) {
		return NewStsCacheProvider(mfaCacheKey, m.InCache)
	}

	if awsConfig.NonChainedGetSessionTokenDuration == 0 {
		awsConfig.NonChainedGetSessionTokenDuration = 15 * time.Minute // default to minimum duration of 15 minutes
	}

	return &mfaSessionTokenProvider{
		SessionTokenProvider: confighelpers.SessionTokenProvider{
			StsClient: getSTSClient(awsConfig.Region, m.NewAccessKeysProvider()),
			Duration:  awsConfig.NonChainedGetSessionTokenDuration,
			Mfa:       confighelpers.NewMfa(awsConfig),
		},
		stsCacheWriter: NewSTSCacheWriter(mfaCacheKey, m.OutCache),
	}
}

func (m CacheProviderFactory) NewAccessKeysProvider() aws.CredentialsProvider {
	return accessKeysProvider{itemFields: m.ItemFields}
}

func (m CacheProviderFactory) NewSSOProvider(awsConfig *confighelpers.Config) aws.CredentialsProvider {
	ssoTempCredCacheKey := getSSOTemporaryCredentialsCacheKey()
	if m.InCache.Has(ssoTempCredCacheKey) {
		return NewStsCacheProvider(ssoTempCredCacheKey, m.InCache)
	}

	cfg := aws.Config{
		Region: awsConfig.SSORegion,
	}

	return ssoProvider{
		OIDCClient:     ssooidc.NewFromConfig(cfg),
		SSOClient:      sso.NewFromConfig(cfg),
		StartURL:       awsConfig.SSOStartURL,
		AccountID:      awsConfig.SSOAccountID,
		RoleName:       awsConfig.SSORoleName,
		stsCacheWriter: NewSTSCacheWriter(ssoTempCredCacheKey, m.OutCache),
	}
}

// getAWSAuthConfigurationForProfile loads specified configurations from both config file and environment
func getAWSAuthConfigurationForProfile(profile string) (*confighelpers.Config, error) {
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
	if awsConfig.AssumeRoleDuration == 0 {
		awsConfig.AssumeRoleDuration = 15 * time.Minute // default to minimum duration of 15 minutes
	}
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

type ssoProvider struct {
	OIDCClient *ssooidc.Client
	SSOClient  *sso.Client
	StartURL   string
	AccountID  string
	RoleName   string
	stsCacheWriter
}

func (p ssoProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	clientCreds, err := p.OIDCClient.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName: aws.String("aws-shell-plugin"),
		ClientType: aws.String("public"),
	})
	if err != nil {
		return aws.Credentials{}, err
	}

	deviceAuthorization, err := p.OIDCClient.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     clientCreds.ClientId,
		ClientSecret: clientCreds.ClientSecret,
		StartUrl:     aws.String(p.StartURL),
	})
	if err != nil {
		return aws.Credentials{}, err
	}

	log.Printf("Opening your Identity Provider in your default browser (use Ctrl-C to abort)\n%s\n", aws.ToString(deviceAuthorization.VerificationUriComplete))
	if err := open.Run(aws.ToString(deviceAuthorization.VerificationUriComplete)); err != nil {
		return aws.Credentials{}, fmt.Errorf("failed to open your default browser: %w", err)
	}

	// Start polling as described in https://tools.ietf.org/html/draft-ietf-oauth-device-flow-15#section-3.5
	grantType := aws.String("urn:ietf:params:oauth:grant-type:device_code")

	var slowDownDelay = 5 * time.Second
	var retryInterval time.Duration
	if i := deviceAuthorization.Interval; i > 0 {
		retryInterval = time.Duration(i) * time.Second
	} else {
		retryInterval = 5 * time.Second
	}

	var token *ssooidc.CreateTokenOutput
	for {
		t, err := p.OIDCClient.CreateToken(ctx, &ssooidc.CreateTokenInput{
			ClientId:     clientCreds.ClientId,
			ClientSecret: clientCreds.ClientSecret,
			DeviceCode:   deviceAuthorization.DeviceCode,
			GrantType:    grantType,
		})
		if err != nil {
			var pollingTooOftenErr *types.SlowDownException
			if errors.As(err, &pollingTooOftenErr) {
				retryInterval += slowDownDelay
			}

			var notYetAuthorizedErr *types.AuthorizationPendingException
			if errors.As(err, &notYetAuthorizedErr) {
				time.Sleep(retryInterval)
				continue
			}

			return aws.Credentials{}, err
		}

		token = t
		break
	}
	err = exec.Command("open", "-j", "-a", "terminal").Start()
	if err != nil {
		return aws.Credentials{}, err
	}

	resp, err := p.SSOClient.GetRoleCredentials(ctx, &sso.GetRoleCredentialsInput{
		AccessToken: token.AccessToken,
		AccountId:   aws.String(p.AccountID),
		RoleName:    aws.String(p.RoleName),
	})

	temporaryCredentials := aws.Credentials{
		AccessKeyID:     aws.ToString(resp.RoleCredentials.AccessKeyId),
		SecretAccessKey: aws.ToString(resp.RoleCredentials.SecretAccessKey),
		SessionToken:    aws.ToString(resp.RoleCredentials.SessionToken),
		Expires:         time.Unix(0, resp.RoleCredentials.Expiration*int64(time.Millisecond)),
	}

	err = p.stsCacheWriter.Put(temporaryCredentials)
	if err != nil {
		return aws.Credentials{}, err
	}

	return temporaryCredentials, nil
}
