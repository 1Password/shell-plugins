package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
)

func TestSSOProvisioner(t *testing.T) {
	t.Run("cache hit legacy form", func(t *testing.T) {
		t.Setenv("AWS_PROFILE", "")
		t.Setenv("AWS_DEFAULT_REGION", "")
		configPath := filepath.Join(t.TempDir(), "awsConfig")
		t.Setenv("AWS_CONFIG_FILE", configPath)

		file := ini.Empty()
		profile, err := file.NewSection("profile sso-legacy")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_start_url", "https://example.awsapps.com/start")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_region", "us-east-1")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_account_id", "111111111111")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_role_name", "ReadOnly")
		require.NoError(t, err)
		_, err = profile.NewKey("region", "us-east-1")
		require.NoError(t, err)
		require.NoError(t, file.SaveTo(configPath))

		cachedCreds := aws.Credentials{
			AccessKeyID:     "AKIACACHEDLEGACYSSO",
			SecretAccessKey: "cachedSecretLegacy/K7MDENG/bPxRfiCYEXAMPLEKEY",
			SessionToken:    "cachedSessionTokenLegacy",
			CanExpire:       true,
			Expires:         time.Now().Add(30 * time.Minute),
		}
		cacheKey := getSSORoleCacheKey("111111111111", "ReadOnly", "https://example.awsapps.com/start")
		marshaled, err := json.Marshal(cachedCreds)
		require.NoError(t, err)

		factory := &mockSSOProviderFactory{}
		provisioner := SSOProvisioner{
			profileName: "sso-legacy",
			newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) SSOProviderFactory {
				factory.inCache = cacheState
				return factory
			},
		}

		in := sdk.ProvisionInput{
			ItemFields: map[sdk.FieldName]string{},
			HomeDir:    "~",
			TempDir:    "/tmp",
			Cache: sdk.CacheState{
				cacheKey: sdk.CacheEntry{Data: marshaled, ExpiresAt: cachedCreds.Expires},
			},
		}
		out := sdk.ProvisionOutput{
			Environment: make(map[string]string),
			Files:       make(map[string]sdk.OutputFile),
		}

		provisioner.Provision(context.Background(), in, &out)

		require.Empty(t, out.Diagnostics.Errors)
		assert.False(t, factory.freshCalled, "fresh SSO retrieval should not run on cache hit")
		assert.Equal(t, cachedCreds.AccessKeyID, out.Environment["AWS_ACCESS_KEY_ID"])
		assert.Equal(t, cachedCreds.SecretAccessKey, out.Environment["AWS_SECRET_ACCESS_KEY"])
		assert.Equal(t, cachedCreds.SessionToken, out.Environment["AWS_SESSION_TOKEN"])
		assert.Equal(t, "us-east-1", out.Environment["AWS_DEFAULT_REGION"])
	})

	t.Run("cache hit sso_session form", func(t *testing.T) {
		t.Setenv("AWS_PROFILE", "")
		t.Setenv("AWS_DEFAULT_REGION", "")
		configPath := filepath.Join(t.TempDir(), "awsConfig")
		t.Setenv("AWS_CONFIG_FILE", configPath)

		file := ini.Empty()
		profile, err := file.NewSection("profile sso-session-prof")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_session", "corp")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_account_id", "222222222222")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_role_name", "PowerUser")
		require.NoError(t, err)
		_, err = profile.NewKey("region", "us-west-2")
		require.NoError(t, err)

		session, err := file.NewSection("sso-session corp")
		require.NoError(t, err)
		_, err = session.NewKey("sso_start_url", "https://corp.awsapps.com/start")
		require.NoError(t, err)
		_, err = session.NewKey("sso_region", "us-west-2")
		require.NoError(t, err)
		require.NoError(t, file.SaveTo(configPath))

		cachedCreds := aws.Credentials{
			AccessKeyID:     "AKIACACHEDSSOSESSION",
			SecretAccessKey: "cachedSecretSession/K7MDENG/bPxRfiCYEXAMPLEKEY",
			SessionToken:    "cachedSessionTokenSession",
			CanExpire:       true,
			Expires:         time.Now().Add(30 * time.Minute),
		}
		cacheKey := getSSORoleCacheKey("222222222222", "PowerUser", "corp")
		marshaled, err := json.Marshal(cachedCreds)
		require.NoError(t, err)

		factory := &mockSSOProviderFactory{}
		provisioner := SSOProvisioner{
			profileName: "sso-session-prof",
			newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) SSOProviderFactory {
				factory.inCache = cacheState
				return factory
			},
		}

		in := sdk.ProvisionInput{
			ItemFields: map[sdk.FieldName]string{},
			HomeDir:    "~",
			TempDir:    "/tmp",
			Cache: sdk.CacheState{
				cacheKey: sdk.CacheEntry{Data: marshaled, ExpiresAt: cachedCreds.Expires},
			},
		}
		out := sdk.ProvisionOutput{
			Environment: make(map[string]string),
			Files:       make(map[string]sdk.OutputFile),
		}

		provisioner.Provision(context.Background(), in, &out)

		require.Empty(t, out.Diagnostics.Errors)
		assert.False(t, factory.freshCalled, "fresh SSO retrieval should not run on cache hit")
		assert.Equal(t, cachedCreds.AccessKeyID, out.Environment["AWS_ACCESS_KEY_ID"])
		assert.Equal(t, cachedCreds.SecretAccessKey, out.Environment["AWS_SECRET_ACCESS_KEY"])
		assert.Equal(t, cachedCreds.SessionToken, out.Environment["AWS_SESSION_TOKEN"])
		assert.Equal(t, "us-west-2", out.Environment["AWS_DEFAULT_REGION"])
	})

	t.Run("cache miss with valid SSO token", func(t *testing.T) {
		t.Setenv("AWS_PROFILE", "")
		t.Setenv("AWS_DEFAULT_REGION", "")
		configPath := filepath.Join(t.TempDir(), "awsConfig")
		t.Setenv("AWS_CONFIG_FILE", configPath)

		file := ini.Empty()
		profile, err := file.NewSection("profile sso-fresh")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_start_url", "https://example.awsapps.com/start")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_region", "us-east-1")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_account_id", "333333333333")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_role_name", "Admin")
		require.NoError(t, err)
		_, err = profile.NewKey("region", "us-east-1")
		require.NoError(t, err)
		require.NoError(t, file.SaveTo(configPath))

		plugintest.TestProvisioner(t, SSOProvisioner{
			profileName: "sso-fresh",
			newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) SSOProviderFactory {
				return &mockSSOProviderFactory{provider: mockSSORoleProvider{}}
			},
		}, map[string]plugintest.ProvisionCase{
			"emits SSO role credentials": {
				ItemFields: map[sdk.FieldName]string{},
				ExpectedOutput: sdk.ProvisionOutput{
					Environment: map[string]string{
						"AWS_ACCESS_KEY_ID":     "AKIASSOFRESHCREDS",
						"AWS_SECRET_ACCESS_KEY": "ssoFreshSecret/K7MDENG/bPxRfiCYEXAMPLEKEY",
						"AWS_SESSION_TOKEN":     "ssoFreshSessionToken/K7MDENG/bPxRfiCYEXAMPLEKEY",
						"AWS_DEFAULT_REGION":    "us-east-1",
					},
				},
			},
		})
	})

	t.Run("cache miss with invalid token error", func(t *testing.T) {
		t.Setenv("AWS_PROFILE", "")
		t.Setenv("AWS_DEFAULT_REGION", "")
		configPath := filepath.Join(t.TempDir(), "awsConfig")
		t.Setenv("AWS_CONFIG_FILE", configPath)

		file := ini.Empty()
		profile, err := file.NewSection("profile sso-expired")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_start_url", "https://example.awsapps.com/start")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_region", "us-east-1")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_account_id", "444444444444")
		require.NoError(t, err)
		_, err = profile.NewKey("sso_role_name", "Auditor")
		require.NoError(t, err)
		require.NoError(t, file.SaveTo(configPath))

		plugintest.TestProvisioner(t, SSOProvisioner{
			profileName: "sso-expired",
			newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) SSOProviderFactory {
				return &mockSSOProviderFactory{provider: mockSSOInvalidTokenProvider{}}
			},
		}, map[string]plugintest.ProvisionCase{
			"surfaces friendly login instruction": {
				ItemFields: map[sdk.FieldName]string{},
				ExpectedOutput: sdk.ProvisionOutput{
					Diagnostics: sdk.Diagnostics{Errors: []sdk.Error{{Message: "AWS SSO token is missing or expired; run `aws sso login --profile sso-expired` and try again"}}},
				},
			},
		})
	})

	t.Run("profile is not SSO", func(t *testing.T) {
		// When the active profile is a plain IAM profile, the SSO provisioner must
		// yield silently — no env vars, no error — so the Access Key provisioner
		// can run.
		t.Setenv("AWS_PROFILE", "")
		t.Setenv("AWS_DEFAULT_REGION", "")
		configPath := filepath.Join(t.TempDir(), "awsConfig")
		t.Setenv("AWS_CONFIG_FILE", configPath)

		file := ini.Empty()
		profile, err := file.NewSection("profile iam-user")
		require.NoError(t, err)
		_, err = profile.NewKey("region", "us-east-1")
		require.NoError(t, err)
		require.NoError(t, file.SaveTo(configPath))

		plugintest.TestProvisioner(t, SSOProvisioner{
			profileName: "iam-user",
			newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) SSOProviderFactory {
				return &mockSSOProviderFactory{}
			},
		}, map[string]plugintest.ProvisionCase{
			"yields silently": {
				ItemFields:     map[sdk.FieldName]string{},
				ExpectedOutput: sdk.ProvisionOutput{},
			},
		})
	})
}

func TestResolveLocalAnd1PasswordSSOConfigurations(t *testing.T) {
	for _, scenario := range []struct {
		description             string
		itemFields              map[sdk.FieldName]string
		awsConfig               *confighelpers.Config
		expectedResultingConfig *confighelpers.Config
		err                     error
	}{
		{
			description: "all SSO fields present only in 1Password",
			itemFields: map[sdk.FieldName]string{
				fieldname.SSOStartURL:  "https://example.awsapps.com/start",
				fieldname.SSORegion:    "us-east-1",
				fieldname.SSOAccountID: "111111111111",
				fieldname.SSORoleName:  "ReadOnly",
				fieldname.SSOSession:   "corp",
			},
			awsConfig: &confighelpers.Config{ProfileName: "dev"},
			expectedResultingConfig: &confighelpers.Config{
				ProfileName:  "dev",
				SSOStartURL:  "https://example.awsapps.com/start",
				SSORegion:    "us-east-1",
				SSOAccountID: "111111111111",
				SSORoleName:  "ReadOnly",
				SSOSession:   "corp",
			},
		},
		{
			description: "all SSO fields present only in local config",
			itemFields:  map[sdk.FieldName]string{},
			awsConfig: &confighelpers.Config{
				ProfileName:  "dev",
				SSOStartURL:  "https://example.awsapps.com/start",
				SSORegion:    "us-east-1",
				SSOAccountID: "111111111111",
				SSORoleName:  "ReadOnly",
			},
			expectedResultingConfig: &confighelpers.Config{
				ProfileName:  "dev",
				SSOStartURL:  "https://example.awsapps.com/start",
				SSORegion:    "us-east-1",
				SSOAccountID: "111111111111",
				SSORoleName:  "ReadOnly",
			},
		},
		{
			description: "SSO fields agree between 1Password and local config",
			itemFields: map[sdk.FieldName]string{
				fieldname.SSOStartURL:  "https://example.awsapps.com/start",
				fieldname.SSORegion:    "us-east-1",
				fieldname.SSOAccountID: "111111111111",
				fieldname.SSORoleName:  "ReadOnly",
			},
			awsConfig: &confighelpers.Config{
				ProfileName:  "dev",
				SSOStartURL:  "https://example.awsapps.com/start",
				SSORegion:    "us-east-1",
				SSOAccountID: "111111111111",
				SSORoleName:  "ReadOnly",
			},
			expectedResultingConfig: &confighelpers.Config{
				ProfileName:  "dev",
				SSOStartURL:  "https://example.awsapps.com/start",
				SSORegion:    "us-east-1",
				SSOAccountID: "111111111111",
				SSORoleName:  "ReadOnly",
			},
		},
		{
			description: "SSO start URL conflict",
			itemFields: map[sdk.FieldName]string{
				fieldname.SSOStartURL: "https://other.awsapps.com/start",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				SSOStartURL: "https://example.awsapps.com/start",
			},
			err: fmt.Errorf("your local AWS configuration has a different SSO Start URL than the one specified in 1Password"),
		},
		{
			description: "SSO region conflict",
			itemFields: map[sdk.FieldName]string{
				fieldname.SSORegion: "us-west-2",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				SSORegion:   "us-east-1",
			},
			err: fmt.Errorf("your local AWS configuration has a different SSO Region than the one specified in 1Password"),
		},
		{
			description: "SSO account ID conflict",
			itemFields: map[sdk.FieldName]string{
				fieldname.SSOAccountID: "999999999999",
			},
			awsConfig: &confighelpers.Config{
				ProfileName:  "dev",
				SSOAccountID: "111111111111",
			},
			err: fmt.Errorf("your local AWS configuration has a different SSO Account ID than the one specified in 1Password"),
		},
		{
			description: "SSO role name conflict",
			itemFields: map[sdk.FieldName]string{
				fieldname.SSORoleName: "Admin",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				SSORoleName: "ReadOnly",
			},
			err: fmt.Errorf("your local AWS configuration has a different SSO Role Name than the one specified in 1Password"),
		},
		{
			description: "SSO session conflict",
			itemFields: map[sdk.FieldName]string{
				fieldname.SSOSession: "personal",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				SSOSession:  "corp",
			},
			err: fmt.Errorf("your local AWS configuration has a different SSO Session than the one specified in 1Password"),
		},
		{
			description: "default region present only in 1Password",
			itemFields: map[sdk.FieldName]string{
				fieldname.DefaultRegion: "us-east-2",
			},
			awsConfig: &confighelpers.Config{
				ProfileName:  "dev",
				SSOStartURL:  "https://example.awsapps.com/start",
				SSOAccountID: "111111111111",
				SSORoleName:  "ReadOnly",
			},
			expectedResultingConfig: &confighelpers.Config{
				ProfileName:  "dev",
				SSOStartURL:  "https://example.awsapps.com/start",
				SSOAccountID: "111111111111",
				SSORoleName:  "ReadOnly",
				Region:       "us-east-2",
			},
		},
		{
			description: "default region conflict between 1Password and local config",
			itemFields: map[sdk.FieldName]string{
				fieldname.DefaultRegion: "us-east-2",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				Region:      "us-east-1",
			},
			err: fmt.Errorf("your local AWS configuration has a different default region than the one specified in 1Password"),
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			err := resolveLocalAnd1PasswordSSOConfigurations(scenario.itemFields, scenario.awsConfig)
			if scenario.err != nil {
				assert.EqualError(t, err, scenario.err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, scenario.expectedResultingConfig, scenario.awsConfig)
			}
		})
	}
}

func TestMissingRequiredSSOFields(t *testing.T) {
	for _, scenario := range []struct {
		description string
		awsConfig   *confighelpers.Config
		expected    []string
	}{
		{
			description: "all fields present",
			awsConfig: &confighelpers.Config{
				SSOStartURL:  "https://example.awsapps.com/start",
				SSORegion:    "us-east-1",
				SSOAccountID: "111111111111",
				SSORoleName:  "ReadOnly",
			},
			expected: nil,
		},
		{
			description: "all fields missing",
			awsConfig:   &confighelpers.Config{},
			expected:    []string{"SSO Start URL", "SSO Region", "SSO Account ID", "SSO Role Name"},
		},
		{
			description: "only SSO Start URL missing (orphan sso_session reference)",
			awsConfig: &confighelpers.Config{
				SSOSession:   "corp",
				SSORegion:    "us-east-1",
				SSOAccountID: "111111111111",
				SSORoleName:  "ReadOnly",
			},
			expected: []string{"SSO Start URL"},
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			assert.Equal(t, scenario.expected, missingRequiredSSOFields(scenario.awsConfig))
		})
	}
}

type mockSSOProviderFactory struct {
	inCache     sdk.CacheState
	provider    aws.CredentialsProvider
	freshCalled bool
}

func (f *mockSSOProviderFactory) NewSSORoleCredentialsProvider(awsConfig *confighelpers.Config) aws.CredentialsProvider {
	cacheKey := getSSORoleCacheKey(awsConfig.SSOAccountID, awsConfig.SSORoleName, ssoSessionKey(awsConfig))
	if f.inCache.Has(cacheKey) {
		return NewStsCacheProvider(cacheKey, f.inCache)
	}
	f.freshCalled = true
	return f.provider
}

type mockSSORoleProvider struct{}

func (mockSSORoleProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "AKIASSOFRESHCREDS",
		SecretAccessKey: "ssoFreshSecret/K7MDENG/bPxRfiCYEXAMPLEKEY",
		SessionToken:    "ssoFreshSessionToken/K7MDENG/bPxRfiCYEXAMPLEKEY",
		CanExpire:       true,
		Expires:         time.Now().Add(time.Hour),
	}, nil
}

type mockSSOInvalidTokenProvider struct{}

func (mockSSOInvalidTokenProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{}, &ssocreds.InvalidTokenError{Err: fmt.Errorf("token cache file does not exist")}
}
