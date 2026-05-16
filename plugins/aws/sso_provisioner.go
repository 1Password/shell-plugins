package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	smithy "github.com/aws/smithy-go"
)

// ssoRetrieveTimeout caps a single sso:GetRoleCredentials round-trip. The SSO endpoint can be
// reached via an attacker-influenceable region if the local AWS config is compromised; a per-call
// deadline ensures one bad config cannot wedge the plugin indefinitely.
const ssoRetrieveTimeout = 30 * time.Second

// SSOProvisioner provisions short-lived AWS credentials by exchanging an SSO access token
// (cached by `aws sso login`) for role credentials via sso:GetRoleCredentials.
type SSOProvisioner struct {
	profileName        string
	newProviderFactory func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) SSOProviderFactory
}

func NewSSOProvisioner(profileName string) SSOProvisioner {
	return SSOProvisioner{
		profileName: profileName,
		newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) SSOProviderFactory {
			return &SSOCacheProviderFactory{
				InCache:    cacheState,
				OutCache:   cacheOps,
				ItemFields: fields,
			}
		},
	}
}

func (p SSOProvisioner) getProfile() (string, error) {
	if len(p.profileName) != 0 {
		return p.profileName, nil
	}

	if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		return profile, nil
	}

	return defaultProfileName, nil
}

func (p SSOProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	profile, err := p.getProfile()
	if err != nil {
		out.AddError(err)
		return
	}

	awsConfig, err := ExecuteSilently(getAWSAuthConfigurationForProfile)(profile)
	if err != nil {
		out.AddError(err)
		return
	}

	// If the selected profile is not configured for AWS IAM Identity Center (SSO),
	// this provisioner is not the right one to handle the request. Yield silently
	// so the Access Key provisioner — which the user has linked in the same `Uses`
	// block — can supply credentials.
	if !awsConfig.HasSSOStartURL() && !awsConfig.HasSSOSession() {
		return
	}

	if err := resolveLocalAnd1PasswordSSOConfigurations(in.ItemFields, awsConfig); err != nil {
		out.AddError(err)
		return
	}

	if missing := missingRequiredSSOFields(awsConfig); len(missing) > 0 {
		out.AddError(fmt.Errorf("missing required field(s) for AWS SSO: %s; add them to the 1Password item or to profile %q in your AWS config file", strings.Join(missing, ", "), profile))
		return
	}

	factory := p.newProviderFactory(in.Cache, out.Cache, in.ItemFields)
	credsProvider := factory.NewSSORoleCredentialsProvider(awsConfig)

	retrieveCtx, cancel := context.WithTimeout(ctx, ssoRetrieveTimeout)
	defer cancel()

	creds, err := ExecuteSilently(credsProvider.Retrieve)(retrieveCtx)
	if err != nil {
		out.AddError(translateSSORetrieveError(err, profile))
		return
	}

	out.AddEnvVar("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
	out.AddEnvVar("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
	if creds.SessionToken != "" {
		out.AddEnvVar("AWS_SESSION_TOKEN", creds.SessionToken)
	}
	if awsConfig.Region != "" {
		out.AddEnvVar("AWS_DEFAULT_REGION", awsConfig.Region)
	}
}

func (p SSOProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p SSOProvisioner) Description() string {
	return "Provision environment variables with temporary AWS IAM Identity Center role credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}

// SSOProviderFactory builds an aws.CredentialsProvider that returns role credentials
// derived from a cached AWS SSO access token.
type SSOProviderFactory interface {
	NewSSORoleCredentialsProvider(awsConfig *confighelpers.Config) aws.CredentialsProvider
}

// SSOCacheProviderFactory wraps the underlying SSO role-credentials provider with the
// shell plugin's encrypted cache so subsequent runs within the credential's TTL avoid
// hitting the SSO endpoint.
type SSOCacheProviderFactory struct {
	InCache    sdk.CacheState
	OutCache   sdk.CacheOperations
	ItemFields map[sdk.FieldName]string
}

func (f SSOCacheProviderFactory) NewSSORoleCredentialsProvider(awsConfig *confighelpers.Config) aws.CredentialsProvider {
	cacheKey := getSSORoleCacheKey(awsConfig.SSOAccountID, awsConfig.SSORoleName, ssoSessionKey(awsConfig))
	if f.InCache.Has(cacheKey) {
		return NewStsCacheProvider(cacheKey, f.InCache)
	}

	cachedTokenFilepath, err := ssocreds.StandardCachedTokenFilepath(ssoSessionKey(awsConfig))
	if err != nil {
		return errProvider{err: err}
	}

	if err := assertSSOTokenCacheSafe(cachedTokenFilepath); err != nil {
		return errProvider{err: err}
	}

	ssoClient := sso.NewFromConfig(aws.Config{Region: awsConfig.SSORegion})

	provider := ssocreds.New(ssoClient, awsConfig.SSOAccountID, awsConfig.SSORoleName, awsConfig.SSOStartURL, func(o *ssocreds.Options) {
		o.CachedTokenFilepath = cachedTokenFilepath
	})

	return &ssoRoleCacheWritingProvider{
		Provider:       provider,
		stsCacheWriter: NewSTSCacheWriter(cacheKey, f.OutCache),
	}
}

// ssoSessionKey returns the value used to derive the shared `~/.aws/sso/cache/<sha1>.json`
// filename. botocore uses the sso_session name when present, otherwise the start URL.
func ssoSessionKey(awsConfig *confighelpers.Config) string {
	if awsConfig.HasSSOSession() {
		return awsConfig.SSOSession
	}
	return awsConfig.SSOStartURL
}

// ssoRoleCacheWritingProvider wraps the SDK's SSO role credentials provider so that
// successful retrievals are persisted in the shell plugin's encrypted cache.
type ssoRoleCacheWritingProvider struct {
	*ssocreds.Provider
	stsCacheWriter
}

func (p *ssoRoleCacheWritingProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	creds, err := p.Provider.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}
	if err := p.stsCacheWriter.Put(creds); err != nil {
		return aws.Credentials{}, err
	}
	return creds, nil
}

// errProvider returns a fixed error from Retrieve. It exists so the cache-key derivation
// path can surface filesystem errors through the same code path as a normal Retrieve call.
type errProvider struct {
	err error
}

func (p errProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{}, p.err
}

// translateSSORetrieveError rewrites token-not-found / token-expired errors from ssocreds into
// a friendly message that points the user at `aws sso login`. For other smithy.APIError variants
// returned by the SSO endpoint, it maps known codes to plugin-controlled strings; unknown codes
// get a generic message so attacker-controlled error text from a hostile endpoint never reaches
// the user-visible UI. Non-smithy errors (e.g. our own filesystem fail-closed errors from
// assertSSOTokenCacheSafe) are passed through unchanged because they are locally produced.
func translateSSORetrieveError(err error, profile string) error {
	var invalid *ssocreds.InvalidTokenError
	if errors.As(err, &invalid) {
		cmd := "aws sso login"
		if profile != defaultProfileName {
			cmd = fmt.Sprintf("aws sso login --profile %s", profile)
		}
		return fmt.Errorf("AWS SSO token is missing or expired; run `%s` and try again", cmd)
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "UnauthorizedException":
			return fmt.Errorf("AWS SSO rejected the cached access token as unauthorized; run `aws sso login` and try again")
		case "ForbiddenException":
			return fmt.Errorf("AWS SSO denied access for the configured account/role; verify the assigned permissions in IAM Identity Center")
		case "ResourceNotFoundException":
			return fmt.Errorf("AWS SSO could not find the configured account or role; verify sso_account_id and sso_role_name")
		case "TooManyRequestsException":
			return fmt.Errorf("AWS SSO is throttling requests; wait a moment and try again")
		default:
			// Internally log the original error so an operator can still diagnose; do not surface
			// the server-controlled message text in the user-facing error.
			log.Printf("aws sso plugin: unexpected SSO API error code %q from sso:GetRoleCredentials", apiErr.ErrorCode())
			return fmt.Errorf("failed to retrieve SSO role credentials; check AWS configuration and try again")
		}
	}

	return err
}

// assertSSOTokenCacheSafe enforces fail-closed properties on `~/.aws/sso/cache/<sha1>.json`
// before the AWS SDK reads it: refuses symlinks (which would let a co-resident attacker
// substitute a forged token), refuses files with group/world-readable bits set (since SSO
// access tokens are bearer credentials), and on Unix refuses files not owned by the current
// user. A non-existent file is allowed: the SDK will surface InvalidTokenError, which our
// caller translates into the friendly "run `aws sso login`" message.
func assertSSOTokenCacheSafe(path string) error {
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("inspecting AWS SSO token cache %q: %w", path, err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("AWS SSO token cache %q is a symlink; refusing to follow", path)
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("AWS SSO token cache %q is not a regular file", path)
	}
	if fi.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("AWS SSO token cache %q is group/world readable (mode %o); chmod 600 it and re-run", path, fi.Mode().Perm())
	}
	if st, ok := fi.Sys().(*syscall.Stat_t); ok {
		if st.Uid != uint32(os.Geteuid()) {
			return fmt.Errorf("AWS SSO token cache %q is not owned by the current user", path)
		}
	}
	return nil
}

// missingRequiredSSOFields reports which SSO fields are still empty after merging the 1Password
// item with the local AWS config. Returning early with a clear list avoids late, opaque failures
// from ssocreds or the AWS SSO API.
func missingRequiredSSOFields(awsConfig *confighelpers.Config) []string {
	var missing []string
	if awsConfig.SSOStartURL == "" {
		missing = append(missing, "SSO Start URL")
	}
	if awsConfig.SSORegion == "" {
		missing = append(missing, "SSO Region")
	}
	if awsConfig.SSOAccountID == "" {
		missing = append(missing, "SSO Account ID")
	}
	if awsConfig.SSORoleName == "" {
		missing = append(missing, "SSO Role Name")
	}
	return missing
}

// resolveLocalAnd1PasswordSSOConfigurations reconciles SSO settings between the 1Password
// item and the local AWS config file using the same rules as the existing Access Key flow:
// values present in only one source are accepted; values present in both must agree.
func resolveLocalAnd1PasswordSSOConfigurations(itemFields map[sdk.FieldName]string, awsConfig *confighelpers.Config) error {
	checks := []struct {
		name      string
		fieldName sdk.FieldName
		target    *string
	}{
		{name: "SSO Start URL", fieldName: fieldname.SSOStartURL, target: &awsConfig.SSOStartURL},
		{name: "SSO Region", fieldName: fieldname.SSORegion, target: &awsConfig.SSORegion},
		{name: "SSO Account ID", fieldName: fieldname.SSOAccountID, target: &awsConfig.SSOAccountID},
		{name: "SSO Role Name", fieldName: fieldname.SSORoleName, target: &awsConfig.SSORoleName},
		{name: "SSO Session", fieldName: fieldname.SSOSession, target: &awsConfig.SSOSession},
	}

	for _, c := range checks {
		itemVal, has := itemFields[c.fieldName]
		if !has || itemVal == "" {
			continue
		}
		if *c.target != "" && *c.target != itemVal {
			return fmt.Errorf("your local AWS configuration has a different %s than the one specified in 1Password", c.name)
		}
		*c.target = itemVal
	}

	region, hasRegularRegion := itemFields[fieldname.Region]
	defaultRegion, hasDefaultRegion := itemFields[fieldname.DefaultRegion]
	if hasDefaultRegion {
		region = defaultRegion
	}
	hasRegion := hasRegularRegion || hasDefaultRegion
	if hasRegion && awsConfig.Region != "" && region != awsConfig.Region {
		return fmt.Errorf("your local AWS configuration has a different default region than the one specified in 1Password")
	}
	if awsConfig.Region == "" {
		awsConfig.Region = region
	}

	return nil
}
