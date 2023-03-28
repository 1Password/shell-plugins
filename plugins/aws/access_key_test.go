package aws

import (
	"fmt"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessKey().Importer, map[string]plugintest.ImportCase{
		"AWS CLI default config file location": {
			Files: map[string]string{
				"~/.aws/credentials": plugintest.LoadFixture(t, "credentials"),
				"~/.aws/config":      plugintest.LoadFixture(t, "config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						fieldname.DefaultRegion:   "eu-central-1",
					},
				},
				{
					NameHint: "user1",
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						fieldname.DefaultRegion:   "us-east-1",
					},
				},
			},
		},
		"AWS CLI custom config file in home dir": {
			Environment: map[string]string{
				"AWS_CONFIG_FILE": "~/.config-custom",
			},
			Files: map[string]string{
				"~/.aws/credentials": plugintest.LoadFixture(t, "credentials"),
				"~/.config-custom":   plugintest.LoadFixture(t, "custom-config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						fieldname.DefaultRegion:   "us-west-1",
					},
				},
				{
					NameHint: "user1",
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						fieldname.DefaultRegion:   "us-west-1",
					},
				},
			},
		},
		"AWS CLI custom config file in root dir": {
			Environment: map[string]string{
				"AWS_CONFIG_FILE": "/.config-custom",
			},
			Files: map[string]string{
				"~/.aws/credentials": plugintest.LoadFixture(t, "credentials"),
				"/.config-custom":    plugintest.LoadFixture(t, "custom-config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						fieldname.DefaultRegion:   "us-west-1",
					},
				},
				{
					NameHint: "user1",
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						fieldname.DefaultRegion:   "us-west-1",
					},
				},
			},
		},
		"AWS CLI custom config file in root dir 2": {
			Environment: map[string]string{
				"AWS_CONFIG_FILE": ".config-custom",
			},
			Files: map[string]string{
				"~/.aws/credentials": plugintest.LoadFixture(t, "credentials"),
				".config-custom":     plugintest.LoadFixture(t, "custom-config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						fieldname.DefaultRegion:   "us-west-1",
					},
				},
				{
					NameHint: "user1",
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						fieldname.DefaultRegion:   "us-west-1",
					},
				},
			},
		},
		"AWS CLI empty config file": {
			Files: map[string]string{
				"~/.aws/credentials": plugintest.LoadFixture(t, "credentials"),
				"~/.aws/config":      plugintest.LoadFixture(t, "empty-config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					},
				},
				{
					NameHint: "user1",
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					},
				},
			},
		},
		"AWS CLI NO config file": {
			Files: map[string]string{
				"~/.aws/credentials": plugintest.LoadFixture(t, "credentials"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					},
				},
				{
					NameHint: "user1",
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					},
				},
			},
		},
		"default env vars": {
			Environment: map[string]string{
				"AWS_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXAMPLE",
				"AWS_SECRET_ACCESS_KEY": "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
				"AWS_DEFAULT_REGION":    "us-central-1",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
						fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
						fieldname.DefaultRegion:   "us-central-1",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.DefaultRegion: "us-central-1",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.DefaultRegion: "us-central-1",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.DefaultRegion: "us-central-1",
					},
				},
			},
		},
		"env vars with AMAZON_ prefix": {
			Environment: map[string]string{
				"AMAZON_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXAMPLE",
				"AMAZON_SECRET_ACCESS_KEY": "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
						fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
					},
				},
			},
		},
		"AWS_SECRET_KEY": {
			Environment: map[string]string{
				"AWS_ACCESS_KEY": "AKIAHPIZFMD5EEXAMPLE",
				"AWS_SECRET_KEY": "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
						fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID: "AKIAHPIZFMD5EEXAMPLE",
					},
				},
			},
		},
		"AWS_ACCESS_SECRET": {
			Environment: map[string]string{
				"AWS_ACCESS_KEY":    "AKIAHPIZFMD5EEXAMPLE",
				"AWS_ACCESS_SECRET": "RnnHD6qgcZ0OpYB3chaK73TcobH1YY7yEEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID: "AKIAHPIZFMD5EEXAMPLE",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
						fieldname.SecretAccessKey: "RnnHD6qgcZ0OpYB3chaK73TcobH1YY7yEEXAMPLE",
					},
				},
			},
		},
	})
}

func TestAccessKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
				fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
				fieldname.DefaultRegion:   "us-central-1",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AWS_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXAMPLE",
					"AWS_SECRET_ACCESS_KEY": "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
					"AWS_DEFAULT_REGION":    "us-central-1",
				},
			},
		},
	})
}

func TestResolveLocalAnd1PasswordConfigurations(t *testing.T) {
	for _, scenario := range []struct {
		description             string
		itemFields              map[sdk.FieldName]string
		awsConfig               *confighelpers.Config
		expectedResultingConfig *confighelpers.Config
		err                     error
	}{
		{
			description: "mfa token is already present in local config",
			itemFields: map[sdk.FieldName]string{
				fieldname.OneTimePassword: "000000",
			},
			awsConfig: &confighelpers.Config{
				MfaToken:   "019879",
				MfaProcess: "mfa login",
			},
			err: fmt.Errorf("only 1Password-backed OTP authentication is supported by the MFA worklfow of the AWS shell plugin"),
		},
		{
			description: "mfa data is present only in 1Password",
			itemFields: map[sdk.FieldName]string{
				fieldname.OneTimePassword: "515467",
				fieldname.MFASerial:       "arn:aws:iam::123456789012:mfa/user",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				Region:      "us-east-2",
			},
			expectedResultingConfig: &confighelpers.Config{
				ProfileName: "dev",
				MfaSerial:   "arn:aws:iam::123456789012:mfa/user",
				MfaToken:    "515467",
				Region:      "us-east-2",
			},
		},
		{
			description: "mfa otp is present only in 1Password, while mfa serial is present only in local config",
			itemFields: map[sdk.FieldName]string{
				fieldname.OneTimePassword: "515467",
				fieldname.Region:          "us-east-2",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				MfaSerial:   "arn:aws:iam::123456789012:mfa/user",
				Region:      "us-east-2",
			},
			expectedResultingConfig: &confighelpers.Config{
				ProfileName: "dev",
				MfaSerial:   "arn:aws:iam::123456789012:mfa/user",
				MfaToken:    "515467",
				Region:      "us-east-2",
			},
		},
		{
			description: "mfa serial is present both in 1Password and local config, but their values differ",
			itemFields: map[sdk.FieldName]string{
				fieldname.OneTimePassword: "515467",
				fieldname.MFASerial:       "arn:aws:iam::123456789012:mfa/user1",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				MfaSerial:   "arn:aws:iam::123456789012:mfa/user",
			},
			err: fmt.Errorf("your local AWS configuration (config file or environment variable) has a different MFA serial than the one specified in 1Password"),
		},
		{
			description: "has mfa serial but no mfa token",
			itemFields: map[sdk.FieldName]string{
				fieldname.MFASerial: "arn:aws:iam::123456789012:mfa/user",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				MfaSerial:   "arn:aws:iam::123456789012:mfa/user",
			},
			err: fmt.Errorf("MFA failed: an MFA serial was found but no OTP has been set up in 1Password"),
		},
		{
			description: "has mfa token but no mfa serial",
			itemFields: map[sdk.FieldName]string{
				fieldname.OneTimePassword: "515467",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
			},
			err: fmt.Errorf("MFA failed: an OTP was found wihtout a corresponding MFA serial"),
		},
		{
			description: "has region only in 1Password",
			itemFields: map[sdk.FieldName]string{
				fieldname.OneTimePassword: "515467",
				fieldname.DefaultRegion:   "us-east-2",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				MfaSerial:   "arn:aws:iam::123456789012:mfa/user",
			},
			expectedResultingConfig: &confighelpers.Config{
				ProfileName: "dev",
				MfaSerial:   "arn:aws:iam::123456789012:mfa/user",
				MfaToken:    "515467",
				Region:      "us-east-2",
			},
		},
		{
			description: "has region both in 1Password and local config, but values differ",
			itemFields: map[sdk.FieldName]string{
				fieldname.OneTimePassword: "515467",
				fieldname.DefaultRegion:   "us-east-2",
			},
			awsConfig: &confighelpers.Config{
				ProfileName: "dev",
				MfaSerial:   "arn:aws:iam::123456789012:mfa/user",
				Region:      "us-east-1",
			},
			err: fmt.Errorf("your local AWS configuration (config file or environment variable) has a different default region than the one specified in 1Password"),
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			err := resolveLocalAnd1PasswordConfigurations(scenario.itemFields, scenario.awsConfig)
			if scenario.err != nil {
				assert.EqualError(t, err, scenario.err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, scenario.expectedResultingConfig, scenario.awsConfig)
			}
		})
	}
}

func TestStripAndReturnProfileFlag(t *testing.T) {
	for i, scenario := range []struct {
		args              []string
		expectedArgs      []string
		expectedFlagValue string
		expectedErr       error
	}{
		{
			args:              []string{"--cache", "false", "--session", "asdefg345reger", "--profile", "andy"},
			expectedArgs:      []string{"--cache", "false", "--session", "asdefg345reger"},
			expectedFlagValue: "andy",
		},
		{
			args:              []string{"--cache", "false", "--profile", "andy", "--session", "asdefg345reger"},
			expectedArgs:      []string{"--cache", "false", "--session", "asdefg345reger"},
			expectedFlagValue: "andy",
		},
		{
			args:              []string{"--cache", "false", "--session", "asdefg345reger"},
			expectedArgs:      []string{"--cache", "false", "--session", "asdefg345reger"},
			expectedFlagValue: "",
		},
		{
			args:         []string{"--cache", "false", "--session", "asdefg345reger", "--profile"},
			expectedArgs: []string{"--cache", "false", "--session", "asdefg345reger"},
			expectedErr:  fmt.Errorf("--profile flag was specified without a value"),
		},
		{
			args:         []string{"--cache", "false", "--session", "asdefg345reger", "--profile", "andy", "--profile"},
			expectedArgs: []string{"--cache", "false", "--session", "asdefg345reger"},
			expectedErr:  fmt.Errorf("--profile flag was specified without a value"),
		},
		{
			args:              []string{"--cache", "false", "--session", "asdefg345reger", "--profile=andy"},
			expectedArgs:      []string{"--cache", "false", "--session", "asdefg345reger"},
			expectedFlagValue: "andy",
		},
		{
			args:         []string{"--cache", "false", "--session", "asdefg345reger", "--profile="},
			expectedArgs: []string{"--cache", "false", "--session", "asdefg345reger"},
			expectedErr:  fmt.Errorf("--profile flag was specified without a value"),
		},
		{
			args:              []string{"--profile=andy", "--cache", "false", "--session", "asdefg345reger"},
			expectedArgs:      []string{"--cache", "false", "--session", "asdefg345reger"},
			expectedFlagValue: "andy",
		},
		{
			args:              []string{"--profile=andy", "--cache", "false", "--profile=dev", "--session", "asdefg345reger", "--profile", "prod", "--profile", "stage"},
			expectedArgs:      []string{"--cache", "false", "--session", "asdefg345reger"},
			expectedFlagValue: "stage",
		},
		{
			args:              []string{"--profile=andy", "--profile=dev", "--profile", "prod", "--profile", "stage", "--profile=production"},
			expectedArgs:      nil,
			expectedFlagValue: "production",
		},
		{
			args:              []string{"--profile=andy", "--region=us-east-1", "--", "wrapped command"},
			expectedArgs:      []string{"--region=us-east-1", "--", "wrapped command"},
			expectedFlagValue: "andy",
		},
		{
			args:              []string{"--profile=andy", "--", "wrapped command"},
			expectedArgs:      []string{"--", "wrapped command"},
			expectedFlagValue: "andy",
		},
		{
			args:              []string{"--profile=andy", "--", "wrapped command", "--profile", "dev"},
			expectedArgs:      []string{"--", "wrapped command", "--profile", "dev"},
			expectedFlagValue: "andy",
		},
		{
			args:              []string{"--", "wrapped command", "--profile", "dev"},
			expectedArgs:      []string{"--", "wrapped command", "--profile", "dev"},
			expectedFlagValue: "",
		},
	} {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			flag, args, err := stripAndReturnProfileFlag(scenario.args)
			if scenario.expectedErr != nil {
				assert.EqualError(t, scenario.expectedErr, err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, scenario.expectedFlagValue, flag)
				assert.Equal(t, scenario.expectedArgs, args)
			}
		})

	}
}

func TestGetProfile(t *testing.T) {
	err := os.Setenv("AWS_PROFILE", "dev")
	require.NoError(t, err)

	stsProvisionerWithProfile := STSProvisioner{profileName: "prod"}
	profile, err := stsProvisionerWithProfile.getProfile()
	require.NoError(t, err)
	assert.Equal(t, "prod", profile)

	stsProvisioner := STSProvisioner{}
	profile, err = stsProvisioner.getProfile()
	require.NoError(t, err)
	assert.Equal(t, "dev", profile)

	err = os.Unsetenv("AWS_PROFILE")
	require.NoError(t, err)

	profile, err = stsProvisioner.getProfile()
	require.NoError(t, err)
	assert.Equal(t, "default", profile)
}
