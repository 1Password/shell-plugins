package aws

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
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

func TestAccessKeyDefaultProvisioner(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "awsConfig")
	t.Setenv("AWS_CONFIG_FILE", configPath)

	// setup profiles in config file
	file := ini.Empty()
	profileDefault, err := file.NewSection("default")
	require.NoError(t, err)
	_, err = profileDefault.NewKey("region", "us-central-1")
	require.NoError(t, err)

	err = file.SaveTo(configPath)
	require.NoError(t, err)

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

func TestSTSProvisioner(t *testing.T) {
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_DEFAULT_REGION", "")
	configPath := filepath.Join(t.TempDir(), "awsConfig")
	t.Setenv("AWS_CONFIG_FILE", configPath)

	// setup profiles in config file
	file := ini.Empty()
	profileDev, err := file.NewSection("profile dev")
	require.NoError(t, err)
	_, err = profileDev.NewKey("role_arn", "aws:iam::123456789012:role/testRole2")
	require.NoError(t, err)

	profileProd, err := file.NewSection("profile prod")
	require.NoError(t, err)
	_, err = profileProd.NewKey("mfa_serial", "arn:aws:iam::123456789012:mfa/user1")
	require.NoError(t, err)

	profileDefault, err := file.NewSection("default")
	require.NoError(t, err)
	_, err = profileDefault.NewKey("region", "us-central-1")
	require.NoError(t, err)

	profileTest, err := file.NewSection("profile test")
	require.NoError(t, err)
	_, err = profileTest.NewKey("mfa_serial", "arn:aws:iam::123456789012:mfa/user1")
	require.NoError(t, err)
	_, err = profileTest.NewKey("role_arn", "aws:iam::123456789012:role/testRole")
	require.NoError(t, err)

	profileSourceComplex, err := file.NewSection("profile testSourceComplex")
	require.NoError(t, err)
	_, err = profileSourceComplex.NewKey("mfa_serial", "arn:aws:iam::123456789012:mfa/user1")
	require.NoError(t, err)
	_, err = profileSourceComplex.NewKey("role_arn", "aws:iam::123456789012:role/testRole")
	require.NoError(t, err)
	_, err = profileSourceComplex.NewKey("source_profile", "testSourceSimple")
	require.NoError(t, err)

	profileSourceSimple, err := file.NewSection("profile testSourceSimple")
	require.NoError(t, err)
	_, err = profileSourceSimple.NewKey("source_profile", "default")
	require.NoError(t, err)
	err = file.SaveTo(configPath)
	require.NoError(t, err)

	plugintest.TestProvisioner(t, STSProvisioner{
		profileName: "",
		newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) STSProviderFactory {
			return &mockProviderManager{}
		},
	}, map[string]plugintest.ProvisionCase{
		"WithAccessKeysProvider": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
				fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
				fieldname.DefaultRegion:   "us-central-1",
				fieldname.OneTimePassword: "908789",
				fieldname.MFASerial:       "arn:aws:iam::123456789012:mfa/user1",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AWS_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXSTS",
					"AWS_SECRET_ACCESS_KEY": "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_SESSION_TOKEN":     "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_DEFAULT_REGION":    "us-central-1",
				},
			},
		},
	})
	plugintest.TestProvisioner(t, STSProvisioner{
		profileName: "dev",
		newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) STSProviderFactory {
			return &mockProviderManager{}
		},
	}, map[string]plugintest.ProvisionCase{
		"WithAssumeRoleProvider": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
				fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AWS_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXSTS",
					"AWS_SECRET_ACCESS_KEY": "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_SESSION_TOKEN":     "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_DEFAULT_REGION":    "us-central-1",
				},
			},
		},
	})

	plugintest.TestProvisioner(t, STSProvisioner{
		profileName: "prod",
		newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) STSProviderFactory {
			return &mockProviderManager{}
		},
	}, map[string]plugintest.ProvisionCase{
		"WithMFAProvider": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
				fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
				fieldname.OneTimePassword: "908789",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AWS_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXSTS",
					"AWS_SECRET_ACCESS_KEY": "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_SESSION_TOKEN":     "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_DEFAULT_REGION":    "us-central-1",
				},
			},
		},
	})

	plugintest.TestProvisioner(t, STSProvisioner{
		profileName: "test",
		newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) STSProviderFactory {
			return &mockProviderManager{}
		},
	}, map[string]plugintest.ProvisionCase{
		"WithAssumeRoleAndMFAProvider": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
				fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
				fieldname.OneTimePassword: "908789",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AWS_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXSTS",
					"AWS_SECRET_ACCESS_KEY": "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_SESSION_TOKEN":     "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_DEFAULT_REGION":    "us-central-1",
				},
			},
		},
	})

	plugintest.TestProvisioner(t, STSProvisioner{
		profileName: "testSourceSimple",
		newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) STSProviderFactory {
			return &mockProviderManager{ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
				fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
			}}
		},
	}, map[string]plugintest.ProvisionCase{
		"WithSourceProfileSimple": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
				fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
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

	plugintest.TestProvisioner(t, STSProvisioner{
		profileName: "testSourceComplex",
		newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) STSProviderFactory {
			return &mockProviderManager{}
		},
	}, map[string]plugintest.ProvisionCase{
		"WithSourceProfileComplex": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
				fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
				fieldname.DefaultRegion:   "us-central-1",
				fieldname.OneTimePassword: "908789",
				fieldname.MFASerial:       "arn:aws:iam::123456789012:mfa/user1",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AWS_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXSTS",
					"AWS_SECRET_ACCESS_KEY": "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_SESSION_TOKEN":     "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_DEFAULT_REGION":    "us-central-1",
				},
			},
		},
	})
}

func TestSourceProfileLoop(t *testing.T) {
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_DEFAULT_REGION", "")
	configPath := filepath.Join(t.TempDir(), "awsConfig")
	t.Setenv("AWS_CONFIG_FILE", configPath)

	// setup profiles in config file
	file := ini.Empty()
	profileDev, err := file.NewSection("profile dev")
	require.NoError(t, err)
	_, err = profileDev.NewKey("source_profile", "default")
	require.NoError(t, err)

	profileDefault, err := file.NewSection("default")
	require.NoError(t, err)
	_, err = profileDefault.NewKey("source_profile", "prod")
	require.NoError(t, err)

	profileProd, err := file.NewSection("profile prod")
	require.NoError(t, err)
	_, err = profileProd.NewKey("source_profile", "dev")
	require.NoError(t, err)

	profileStaging, err := file.NewSection("profile staging")
	require.NoError(t, err)
	_, err = profileStaging.NewKey("source_profile", "staging")
	require.NoError(t, err)

	err = file.SaveTo(configPath)
	require.NoError(t, err)

	plugintest.TestProvisioner(t, STSProvisioner{
		profileName: "prod",
		newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) STSProviderFactory {
			return &mockProviderManager{}
		},
	}, map[string]plugintest.ProvisionCase{
		"WithEndlessLoop": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
				fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
				fieldname.DefaultRegion:   "us-central-1",
				fieldname.OneTimePassword: "908789",
				fieldname.MFASerial:       "arn:aws:iam::123456789012:mfa/user1",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Diagnostics: sdk.Diagnostics{Errors: []sdk.Error{{Message: "infinite loop in credential configuration detected. Attempting to load from profile \"prod\" which has already been visited"}}},
			},
		},
	})

	plugintest.TestProvisioner(t, STSProvisioner{
		profileName: "staging",
		newProviderFactory: func(cacheState sdk.CacheState, cacheOps sdk.CacheOperations, fields map[sdk.FieldName]string) STSProviderFactory {
			return &mockProviderManager{}
		},
	}, map[string]plugintest.ProvisionCase{
		"WithAcceptedLoop": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     "AKIAHPIZFMD5EEXAMPLE",
				fieldname.SecretAccessKey: "lBfKB7P5ScmpxDeRoFLZvhJbqNGPoV0vIEXAMPLE",
				fieldname.DefaultRegion:   "us-central-1",
				fieldname.OneTimePassword: "908789",
				fieldname.MFASerial:       "arn:aws:iam::123456789012:mfa/user1",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AWS_ACCESS_KEY_ID":     "AKIAHPIZFMD5EEXSTS",
					"AWS_SECRET_ACCESS_KEY": "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"AWS_SESSION_TOKEN":     "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
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
			err: fmt.Errorf("MFA failed: MFA serial \"arn:aws:iam::123456789012:mfa/user\" was detected on the associated item or in the config file for the selected profile, but no 'One-Time Password' field was found.\nLearn how to add an OTP field to your item:\nhttps://developer.1password.com/docs/cli/shell-plugins/aws/#optional-set-up-multi-factor-authentication"),
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

type mockAwsProvider struct {
}

func (p mockAwsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "AKIAHPIZFMD5EEXSTS",
		SecretAccessKey: "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
		SessionToken:    "stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY///////stststststst/K7MDENG/bPxRfiCYEXAMPLEKEY",
		Source:          "",
		CanExpire:       false,
		Expires:         time.Time{},
	}, nil
}

type mockProviderManager struct {
	ItemFields map[sdk.FieldName]string
}

func (m mockProviderManager) NewMFASessionTokenProvider(awsConfig *confighelpers.Config, srcCredProvider aws.CredentialsProvider) aws.CredentialsProvider {
	return mockAwsProvider{}
}

func (m mockProviderManager) NewAssumeRoleProvider(awsConfig *confighelpers.Config, srcCredProvider aws.CredentialsProvider) aws.CredentialsProvider {
	return mockAwsProvider{}
}

func (m mockProviderManager) NewAccessKeysProvider() aws.CredentialsProvider {
	return accessKeysProvider{itemFields: m.ItemFields}
}
