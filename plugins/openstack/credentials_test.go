package openstack

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"password auth - all fields": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.URL:          "https://keystone.example.com:5000/v3",
				fieldname.Username:     "myuser",
				fieldname.Password:     "s3cr3tpassword",
				fieldname.Project:      "myproject",
				fieldname.ProjectID:    "abc123def456abc123def456abc123de",
				fieldname.Region:       "RegionOne",
				fieldUserDomainName:    "Default",
				fieldProjectDomainName: "Default",
				fieldProjectDomainID:   "default",
				fieldInterface:         "internal",
				fieldIdentityAPIVersion: "3",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"OS_AUTH_URL":             "https://keystone.example.com:5000/v3",
					"OS_USERNAME":             "myuser",
					"OS_PASSWORD":             "s3cr3tpassword",
					"OS_PROJECT_NAME":         "myproject",
					"OS_PROJECT_ID":           "abc123def456abc123def456abc123de",
					"OS_REGION_NAME":          "RegionOne",
					"OS_USER_DOMAIN_NAME":     "Default",
					"OS_PROJECT_DOMAIN_NAME":  "Default",
					"OS_PROJECT_DOMAIN_ID":    "default",
					"OS_INTERFACE":            "internal",
					"OS_IDENTITY_API_VERSION": "3",
					"OS_AUTH_TYPE":            "password",
				},
			},
		},
		"password auth - defaults applied": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.URL:      "https://keystone.example.com:5000/v3",
				fieldname.Username: "myuser",
				fieldname.Password: "s3cr3tpassword",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"OS_AUTH_URL":             "https://keystone.example.com:5000/v3",
					"OS_USERNAME":             "myuser",
					"OS_PASSWORD":             "s3cr3tpassword",
					"OS_INTERFACE":            "public",
					"OS_IDENTITY_API_VERSION": "3",
					"OS_AUTH_TYPE":            "password",
				},
			},
		},
		"application credential auth": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.URL:                   "https://keystone.example.com:5000/v3",
				fieldApplicationCredentialID:    "xxxxxxxxxxxxxxx",
				fieldApplicationCredentialSecret: "yyyyy342lhkwdh",
				fieldname.Region:                "RegionOne",
				fieldProjectDomainName:          "Default",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"OS_AUTH_URL":                      "https://keystone.example.com:5000/v3",
					"OS_APPLICATION_CREDENTIAL_ID":     "xxxxxxxxxxxxxxx",
					"OS_APPLICATION_CREDENTIAL_SECRET": "yyyyy342lhkwdh",
					"OS_REGION_NAME":                   "RegionOne",
					"OS_PROJECT_DOMAIN_NAME":           "Default",
					"OS_INTERFACE":                     "public",
					"OS_IDENTITY_API_VERSION":          "3",
					"OS_AUTH_TYPE":                     "v3applicationcredential",
				},
			},
		},
	})
}

func TestCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, Credentials().Importer, map[string]plugintest.ImportCase{
		"environment variables - password auth": {
			Environment: map[string]string{
				"OS_AUTH_URL":             "https://keystone.example.com:5000/v3",
				"OS_USERNAME":             "myuser",
				"OS_PASSWORD":             "s3cr3tpassword",
				"OS_PROJECT_NAME":         "myproject",
				"OS_PROJECT_ID":           "abc123def456abc123def456abc123de",
				"OS_REGION_NAME":          "RegionOne",
				"OS_USER_DOMAIN_NAME":     "Default",
				"OS_PROJECT_DOMAIN_NAME":  "Default",
				"OS_PROJECT_DOMAIN_ID":    "default",
				"OS_INTERFACE":            "public",
				"OS_IDENTITY_API_VERSION": "3",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.URL:           "https://keystone.example.com:5000/v3",
						fieldname.Username:      "myuser",
						fieldname.Password:      "s3cr3tpassword",
						fieldname.Project:       "myproject",
						fieldname.ProjectID:     "abc123def456abc123def456abc123de",
						fieldname.Region:        "RegionOne",
						fieldUserDomainName:     "Default",
						fieldProjectDomainName:  "Default",
						fieldProjectDomainID:    "default",
						fieldInterface:          "public",
						fieldIdentityAPIVersion: "3",
					},
				},
			},
		},
		"environment variables - application credential auth": {
			Environment: map[string]string{
				"OS_AUTH_URL":                      "https://keystone.example.com:5000/v3",
				"OS_APPLICATION_CREDENTIAL_ID":     "xxxxxxxxxxxxxxx",
				"OS_APPLICATION_CREDENTIAL_SECRET": "yyyyy342lhkwdh",
				"OS_REGION_NAME":                   "RegionOne",
				"OS_AUTH_TYPE":                     "v3applicationcredential",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.URL:                   "https://keystone.example.com:5000/v3",
						fieldApplicationCredentialID:    "xxxxxxxxxxxxxxx",
						fieldApplicationCredentialSecret: "yyyyy342lhkwdh",
						fieldname.Region:                "RegionOne",
						fieldAuthType:                   "v3applicationcredential",
					},
				},
			},
		},
		"cloudrc file": {
			Files: map[string]string{
				"~/openrc.sh": plugintest.LoadFixture(t, "openrc.sh"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.URL:           "https://keystone.example.com:5000/v3",
						fieldname.Username:      "myuser",
						fieldname.Password:      "s3cr3tpassword",
						fieldname.Project:       "myproject",
						fieldname.ProjectID:     "abc123def456abc123def456abc123de",
						fieldname.Region:        "RegionOne",
						fieldUserDomainName:     "Default",
						fieldProjectDomainName:  "Default",
						fieldProjectDomainID:    "default",
						fieldInterface:          "public",
						fieldIdentityAPIVersion: "3",
					},
				},
			},
		},
		"clouds.yaml file - password auth": {
			Files: map[string]string{
				"~/.config/openstack/clouds.yaml": plugintest.LoadFixture(t, "clouds.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.URL:           "https://keystone.example.com:5000/v3",
						fieldname.Username:      "myuser",
						fieldname.Password:      "s3cr3tpassword",
						fieldname.Project:       "myproject",
						fieldname.ProjectID:     "abc123def456abc123def456abc123de",
						fieldname.Region:        "RegionOne",
						fieldUserDomainName:     "Default",
						fieldProjectDomainName:  "Default",
						fieldProjectDomainID:    "default",
						fieldInterface:          "public",
						fieldIdentityAPIVersion: "3",
					},
					NameHint: "mycloud",
				},
			},
		},
		"clouds.yaml file - application credential auth": {
			Files: map[string]string{
				"~/.config/openstack/clouds.yaml": plugintest.LoadFixture(t, "clouds-appcred.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.URL:                   "https://keystone.example.com:5000/v3",
						fieldApplicationCredentialID:    "xxxxxxxxxxxxxxx",
						fieldApplicationCredentialSecret: "yyyyy342lhkwdh",
						fieldname.Region:                "RegionOne",
						fieldProjectDomainName:          "Default",
						fieldInterface:                  "public",
						fieldIdentityAPIVersion:         "3",
						fieldAuthType:                   "v3applicationcredential",
					},
					NameHint: "mycloud",
				},
			},
		},
		"OS_CLIENT_CONFIG_FILE env var": {
			Environment: map[string]string{
				"OS_CLIENT_CONFIG_FILE": "~/custom-clouds.yaml",
			},
			Files: map[string]string{
				"~/custom-clouds.yaml": plugintest.LoadFixture(t, "clouds.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.URL:           "https://keystone.example.com:5000/v3",
						fieldname.Username:      "myuser",
						fieldname.Password:      "s3cr3tpassword",
						fieldname.Project:       "myproject",
						fieldname.ProjectID:     "abc123def456abc123def456abc123de",
						fieldname.Region:        "RegionOne",
						fieldUserDomainName:     "Default",
						fieldProjectDomainName:  "Default",
						fieldProjectDomainID:    "default",
						fieldInterface:          "public",
						fieldIdentityAPIVersion: "3",
					},
					NameHint: "mycloud",
				},
			},
		},
		"clouds.yaml + secure.yaml merged": {
			Files: map[string]string{
				"~/.config/openstack/clouds.yaml": plugintest.LoadFixture(t, "clouds-no-password.yaml"),
				"~/.config/openstack/secure.yaml": plugintest.LoadFixture(t, "secure.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.URL:           "https://keystone.example.com:5000/v3",
						fieldname.Username:      "myuser",
						fieldname.Password:      "s3cr3tpassword",
						fieldname.Project:       "myproject",
						fieldname.ProjectID:     "abc123def456abc123def456abc123de",
						fieldname.Region:        "RegionOne",
						fieldUserDomainName:     "Default",
						fieldProjectDomainName:  "Default",
						fieldProjectDomainID:    "default",
						fieldInterface:          "public",
						fieldIdentityAPIVersion: "3",
					},
					NameHint: "mycloud",
				},
			},
		},
	})
}
