package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"testing"
)

func TestAccessKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessKey().Importer, map[string]plugintest.ImportCase{
		"AWS CLI credentials file": {
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
	})
}
