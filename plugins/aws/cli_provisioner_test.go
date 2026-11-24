package aws

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const cliProvisionerTestConfig = `[default]
region = us-east-1

[profile keys]
region = eu-west-1

[profile sso-legacy]
sso_start_url = https://example.awsapps.com/start
sso_region = us-east-1
sso_account_id = 111111111111
sso_role_name = ReadOnly
region = us-west-2

[sso-session corp]
sso_start_url = https://corp.awsapps.com/start
sso_region = eu-west-1

[profile sso-new]
sso_session = corp
sso_account_id = 222222222222
sso_role_name = PowerUser

[profile chained]
role_arn = arn:aws:iam::111111111111:role/target
source_profile = sso-legacy
`

// TestCLIProvisionersKeepProfileWhenYielding covers the `--profile` handling shared by the Access
// Key and SSO Profile usages: the flag is always stripped, and when the provisioner supplies no
// credentials the selected profile is passed on as AWS_PROFILE instead of being dropped.
func TestCLIProvisionersKeepProfileWhenYielding(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "awsConfig")
	require.NoError(t, os.WriteFile(configPath, []byte(cliProvisionerTestConfig), 0o600))

	accessKeyFields := map[sdk.FieldName]string{
		fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
		fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
	}

	ssoCreds := aws.Credentials{
		AccessKeyID:     "ASIACACHEDSSOEXAMPLE",
		SecretAccessKey: "cachedSecret/K7MDENG/bPxRfiCYEXAMPLEKEY",
		SessionToken:    "cachedSessionToken",
		CanExpire:       true,
		Expires:         time.Now().Add(30 * time.Minute),
	}
	marshaled, err := json.Marshal(ssoCreds)
	require.NoError(t, err)
	ssoCache := sdk.CacheState{
		getSSORoleCacheKey("111111111111", "ReadOnly", "https://example.awsapps.com/start"): {Data: marshaled, ExpiresAt: ssoCreds.Expires},
	}

	cases := map[string]struct {
		provisioner     sdk.Provisioner
		itemFields      map[sdk.FieldName]string
		cache           sdk.CacheState
		ssoLogin        string // session name or start URL of an SSO login to put in the token cache
		envProfile      string
		commandLine     []string
		expectedCommand []string
		expectedEnv     map[string]string
	}{
		"access key yields on a legacy SSO profile and keeps the profile": {
			provisioner:     CLIProvisioner{},
			itemFields:      accessKeyFields,
			commandLine:     []string{"aws", "--profile", "sso-legacy", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv:     map[string]string{"AWS_PROFILE": "sso-legacy"},
		},
		"access key yields on an sso_session profile given as --profile=": {
			provisioner:     CLIProvisioner{},
			itemFields:      accessKeyFields,
			commandLine:     []string{"aws", "--profile=sso-new", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv:     map[string]string{"AWS_PROFILE": "sso-new"},
		},
		"access key provisions for a non-SSO profile and does not set AWS_PROFILE": {
			provisioner:     CLIProvisioner{},
			itemFields:      accessKeyFields,
			commandLine:     []string{"aws", "--profile", "keys", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv: map[string]string{
				"AWS_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXAMPLE",
				"AWS_SECRET_ACCESS_KEY": "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
				"AWS_DEFAULT_REGION":    "eu-west-1",
			},
		},
		"access key yielding on a profile from the AWS_PROFILE variable adds nothing": {
			provisioner:     CLIProvisioner{},
			itemFields:      accessKeyFields,
			envProfile:      "sso-legacy",
			commandLine:     []string{"aws", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv:     map[string]string{},
		},
		"SSO profile yields on a non-SSO profile and keeps the profile": {
			provisioner:     SSOCLIProvisioner{},
			commandLine:     []string{"aws", "--profile", "keys", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv:     map[string]string{"AWS_PROFILE": "keys"},
		},
		"SSO profile yields on a role chain with an SSO source profile and keeps the profile": {
			provisioner:     SSOCLIProvisioner{},
			commandLine:     []string{"aws", "--profile", "chained", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv:     map[string]string{"AWS_PROFILE": "chained"},
		},
		"SSO profile provisions for an SSO profile and does not set AWS_PROFILE": {
			provisioner:     SSOCLIProvisioner{},
			cache:           ssoCache,
			ssoLogin:        "https://example.awsapps.com/start",
			commandLine:     []string{"aws", "--profile", "sso-legacy", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv: map[string]string{
				"AWS_ACCESS_KEY_ID":     "ASIACACHEDSSOEXAMPLE",
				"AWS_SECRET_ACCESS_KEY": "cachedSecret/K7MDENG/bPxRfiCYEXAMPLEKEY",
				"AWS_SESSION_TOKEN":     "cachedSessionToken",
				"AWS_DEFAULT_REGION":    "us-west-2",
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("AWS_CONFIG_FILE", configPath)
			t.Setenv("AWS_PROFILE", c.envProfile)
			t.Setenv("AWS_DEFAULT_REGION", "")
			if c.ssoLogin != "" {
				// Cached role credentials are only used while the SSO login is still there.
				writeSSOTokenCache(t, c.ssoLogin, map[string]string{
					"accessToken": "access-token",
					"expiresAt":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
				})
			}

			itemFields := c.itemFields
			if itemFields == nil {
				itemFields = map[sdk.FieldName]string{}
			}
			in := sdk.ProvisionInput{
				HomeDir:    "~",
				TempDir:    "/tmp",
				ItemFields: itemFields,
				Cache:      c.cache,
			}
			out := sdk.ProvisionOutput{
				Environment: make(map[string]string),
				Files:       make(map[string]sdk.OutputFile),
				CommandLine: c.commandLine,
			}

			c.provisioner.Provision(context.Background(), in, &out)

			require.Empty(t, out.Diagnostics.Errors)
			assert.Equal(t, c.expectedCommand, out.CommandLine)
			assert.Equal(t, c.expectedEnv, out.Environment)
		})
	}
}

func TestKeepProfileWhenYieldingSkipsOnError(t *testing.T) {
	out := sdk.ProvisionOutput{Environment: make(map[string]string)}
	out.AddError(assert.AnError)

	keepProfileWhenYielding("some-profile", &out)

	assert.Empty(t, out.Environment)
}

// In this config the default profile is an SSO profile. The other profiles must not pick up its SSO
// settings: the AWS CLI does not apply the default profile's settings to other profiles either.
const defaultSSOProfileTestConfig = `[default]
sso_start_url = https://example.awsapps.com/start
sso_region = us-east-1
sso_account_id = 111111111111
sso_role_name = ReadOnly
region = us-west-2

[profile keys]
region = eu-west-1

[sso-session corp]
sso_start_url = https://corp.awsapps.com/start
sso_region = eu-west-1

[profile sso-new]
sso_session = corp
sso_account_id = 222222222222
sso_role_name = PowerUser

[profile chained]
role_arn = arn:aws:iam::111111111111:role/target
source_profile = keys
`

func TestProfilesDoNotInheritSSOSettingsFromDefaultProfile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "awsConfig")
	require.NoError(t, os.WriteFile(configPath, []byte(defaultSSOProfileTestConfig), 0o600))
	t.Setenv("AWS_CONFIG_FILE", configPath)
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_DEFAULT_REGION", "")

	usesSSO := func(c *confighelpers.Config) bool {
		return c.HasSSOStartURL() || c.HasSSOSession()
	}

	t.Run("default profile keeps its own SSO settings", func(t *testing.T) {
		awsConfig, err := getAWSAuthConfigurationForProfile("default")
		require.NoError(t, err)
		assert.Equal(t, "https://example.awsapps.com/start", awsConfig.SSOStartURL)
		assert.Equal(t, "111111111111", awsConfig.SSOAccountID)
		assert.Equal(t, "ReadOnly", awsConfig.SSORoleName)
	})

	t.Run("profile without SSO settings is not an SSO profile", func(t *testing.T) {
		awsConfig, err := getAWSAuthConfigurationForProfile("keys")
		require.NoError(t, err)
		assert.False(t, usesSSO(awsConfig))
		assert.Empty(t, awsConfig.SSORegion)
		assert.Empty(t, awsConfig.SSOAccountID)
		assert.Empty(t, awsConfig.SSORoleName)
		assert.Equal(t, "eu-west-1", awsConfig.Region)
	})

	t.Run("profile with its own SSO settings keeps them", func(t *testing.T) {
		awsConfig, err := getAWSAuthConfigurationForProfile("sso-new")
		require.NoError(t, err)
		assert.Equal(t, "corp", awsConfig.SSOSession)
		assert.Equal(t, "https://corp.awsapps.com/start", awsConfig.SSOStartURL)
		assert.Equal(t, "eu-west-1", awsConfig.SSORegion)
		assert.Equal(t, "222222222222", awsConfig.SSOAccountID)
		assert.Equal(t, "PowerUser", awsConfig.SSORoleName)
	})

	t.Run("role chain and its source profile are not SSO profiles", func(t *testing.T) {
		awsConfig, err := getAWSAuthConfigurationForProfile("chained")
		require.NoError(t, err)
		assert.False(t, usesSSO(awsConfig))
		require.NotNil(t, awsConfig.SourceProfile)
		assert.False(t, usesSSO(awsConfig.SourceProfile))
	})
}

func TestCLIProvisionersWithSSODefaultProfile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "awsConfig")
	require.NoError(t, os.WriteFile(configPath, []byte(defaultSSOProfileTestConfig), 0o600))

	accessKeyFields := map[sdk.FieldName]string{
		fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
		fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
	}

	ssoCreds := aws.Credentials{
		AccessKeyID:     "ASIACACHEDSSOEXAMPLE",
		SecretAccessKey: "cachedSecret/K7MDENG/bPxRfiCYEXAMPLEKEY",
		SessionToken:    "cachedSessionToken",
		CanExpire:       true,
		Expires:         time.Now().Add(30 * time.Minute),
	}
	marshaled, err := json.Marshal(ssoCreds)
	require.NoError(t, err)
	ssoCache := sdk.CacheState{
		getSSORoleCacheKey("111111111111", "ReadOnly", "https://example.awsapps.com/start"): {Data: marshaled, ExpiresAt: ssoCreds.Expires},
	}

	cases := map[string]struct {
		provisioner     sdk.Provisioner
		itemFields      map[sdk.FieldName]string
		cache           sdk.CacheState
		ssoLogin        string // session name or start URL of an SSO login to put in the token cache
		commandLine     []string
		expectedCommand []string
		expectedEnv     map[string]string
	}{
		"access key provisions for a profile without SSO settings": {
			provisioner:     CLIProvisioner{},
			itemFields:      accessKeyFields,
			commandLine:     []string{"aws", "--profile", "keys", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv: map[string]string{
				"AWS_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXAMPLE",
				"AWS_SECRET_ACCESS_KEY": "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
				"AWS_DEFAULT_REGION":    "eu-west-1",
			},
		},
		"SSO profile yields on a profile without SSO settings and keeps the profile": {
			provisioner:     SSOCLIProvisioner{},
			cache:           ssoCache,
			ssoLogin:        "https://example.awsapps.com/start",
			commandLine:     []string{"aws", "--profile", "keys", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv:     map[string]string{"AWS_PROFILE": "keys"},
		},
		"access key yields on the SSO default profile": {
			provisioner:     CLIProvisioner{},
			itemFields:      accessKeyFields,
			commandLine:     []string{"aws", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv:     map[string]string{},
		},
		"SSO profile provisions for the SSO default profile": {
			provisioner:     SSOCLIProvisioner{},
			cache:           ssoCache,
			ssoLogin:        "https://example.awsapps.com/start",
			commandLine:     []string{"aws", "s3", "ls"},
			expectedCommand: []string{"aws", "s3", "ls"},
			expectedEnv: map[string]string{
				"AWS_ACCESS_KEY_ID":     "ASIACACHEDSSOEXAMPLE",
				"AWS_SECRET_ACCESS_KEY": "cachedSecret/K7MDENG/bPxRfiCYEXAMPLEKEY",
				"AWS_SESSION_TOKEN":     "cachedSessionToken",
				"AWS_DEFAULT_REGION":    "us-west-2",
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("AWS_CONFIG_FILE", configPath)
			t.Setenv("AWS_PROFILE", "")
			t.Setenv("AWS_DEFAULT_REGION", "")
			if c.ssoLogin != "" {
				writeSSOTokenCache(t, c.ssoLogin, map[string]string{
					"accessToken": "access-token",
					"expiresAt":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
				})
			}

			itemFields := c.itemFields
			if itemFields == nil {
				itemFields = map[sdk.FieldName]string{}
			}
			in := sdk.ProvisionInput{
				HomeDir:    "~",
				TempDir:    "/tmp",
				ItemFields: itemFields,
				Cache:      c.cache,
			}
			out := sdk.ProvisionOutput{
				Environment: make(map[string]string),
				Files:       make(map[string]sdk.OutputFile),
				CommandLine: c.commandLine,
			}

			c.provisioner.Provision(context.Background(), in, &out)

			require.Empty(t, out.Diagnostics.Errors)
			assert.Equal(t, c.expectedCommand, out.CommandLine)
			assert.Equal(t, c.expectedEnv, out.Environment)
		})
	}
}
