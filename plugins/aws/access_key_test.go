package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"testing"
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
					NameHint: "default",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						FieldNameDefaultRegion:    "eu-central-1",
					},
				},
				{
					NameHint: "user1",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						FieldNameDefaultRegion:    "us-east-1",
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
					NameHint: "default",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						FieldNameDefaultRegion:    "us-west-1",
					},
				},
				{
					NameHint: "user1",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						FieldNameDefaultRegion:    "us-west-1",
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
					NameHint: "default",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						FieldNameDefaultRegion:    "us-west-1",
					},
				},
				{
					NameHint: "user1",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						FieldNameDefaultRegion:    "us-west-1",
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
					NameHint: "default",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						FieldNameDefaultRegion:    "us-west-1",
					},
				},
				{
					NameHint: "user1",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						FieldNameDefaultRegion:    "us-west-1",
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
					NameHint: "default",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					},
				},
				{
					NameHint: "user1",
					Fields: map[string]string{
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
					NameHint: "default",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIADEFFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "DEFlrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					},
				},
				{
					NameHint: "user1",
					Fields: map[string]string{
						fieldname.AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
						fieldname.SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					},
				},
			},
		},
	})
}
