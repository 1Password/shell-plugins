package oci

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.User:       "ocid1.user.oc1..example",
				tenantIDField:        "ocid1.tenancy.oc1..example",
				fingerprintField:     "20:3b:97:13:55:1c:example",
				fieldname.Region:     "us-ashburn-1",
				fieldname.PrivateKey: "-----BEGIN PRIVATE KEY-----\nEXAMPLE\n-----END PRIVATE KEY-----",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"OCI_CLI_USER":        "ocid1.user.oc1..example",
					"OCI_CLI_TENANCY":     "ocid1.tenancy.oc1..example",
					"OCI_CLI_FINGERPRINT": "20:3b:97:13:55:1c:example",
					"OCI_CLI_REGION":      "us-ashburn-1",
					"OCI_CLI_KEY_CONTENT": "-----BEGIN PRIVATE KEY-----\nEXAMPLE\n-----END PRIVATE KEY-----",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"OCI_CLI_USER":        "ocid1.user.oc1..example",
				"OCI_CLI_TENANCY":     "ocid1.tenancy.oc1..example",
				"OCI_CLI_FINGERPRINT": "20:3b:97:13:55:1c:example",
				"OCI_CLI_REGION":      "us-ashburn-1",
				"OCI_CLI_KEY_CONTENT": "-----BEGIN PRIVATE KEY-----\nEXAMPLE\n-----END PRIVATE KEY-----",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.User:       "ocid1.user.oc1..example",
						tenantIDField:        "ocid1.tenancy.oc1..example",
						fingerprintField:     "20:3b:97:13:55:1c:example",
						fieldname.Region:     "us-ashburn-1",
						fieldname.PrivateKey: "-----BEGIN PRIVATE KEY-----\nEXAMPLE\n-----END PRIVATE KEY-----",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.oci/config":      plugintest.LoadFixture(t, "config"),
				"~/.oci/default.pem": "-----BEGIN PRIVATE KEY-----\nDEFAULT EXAMPLE\n-----END PRIVATE KEY-----",
				"~/.oci/team.pem":    "-----BEGIN PRIVATE KEY-----\nTEAM EXAMPLE\n-----END PRIVATE KEY-----",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "DEFAULT",
					Fields: map[sdk.FieldName]string{
						fieldname.User:       "ocid1.user.oc1..default-example",
						tenantIDField:        "ocid1.tenancy.oc1..default-example",
						fingerprintField:     "20:3b:97:13:55:1c:default-example",
						fieldname.Region:     "us-ashburn-1",
						fieldname.PrivateKey: "-----BEGIN PRIVATE KEY-----\nDEFAULT EXAMPLE\n-----END PRIVATE KEY-----",
					},
				},
				{
					NameHint: "TEAM",
					Fields: map[sdk.FieldName]string{
						fieldname.User:       "ocid1.user.oc1..team-example",
						tenantIDField:        "ocid1.tenancy.oc1..team-example",
						fingerprintField:     "20:3b:97:13:55:1c:team-example",
						fieldname.PrivateKey: "-----BEGIN PRIVATE KEY-----\nTEAM EXAMPLE\n-----END PRIVATE KEY-----",
					},
				},
			},
		},
	})
}
