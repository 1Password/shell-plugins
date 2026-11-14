package snowflakecli

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const (
	examplePrivateKey          = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQEXAMPLE\n-----END PRIVATE KEY-----"
	exampleEncryptedPrivateKey = "-----BEGIN ENCRYPTED PRIVATE KEY-----\nMIIFHDBOBgkqhkiG9w0BBQ0wQTApBgkqhkiG9w0BBQwwHAQIEXAMPLE\n-----END ENCRYPTED PRIVATE KEY-----"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"required fields only": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Account:    "myorg-myaccount",
				fieldname.Username:   "jdoe",
				fieldname.PrivateKey: examplePrivateKey,
			},
			CommandLine: []string{"snow", "sql", "-q", "select 1"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SNOWFLAKE_ACCOUNT":         "myorg-myaccount",
					"SNOWFLAKE_USER":            "jdoe",
					"SNOWFLAKE_AUTHENTICATOR":   "SNOWFLAKE_JWT",
					"SNOWFLAKE_PRIVATE_KEY_RAW": examplePrivateKey + "\n",
				},
				CommandLine: []string{"snow", "sql", "-q", "select 1", "-x"},
			},
		},
		"all fields": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Account:         "myorg-myaccount",
				fieldname.Username:        "jdoe",
				fieldname.PrivateKey:      exampleEncryptedPrivateKey,
				fieldPrivateKeyPassphrase: "passphrase1234",
				fieldRole:                 "SYSADMIN",
				fieldWarehouse:            "COMPUTE_WH",
				fieldname.Database:        "MY_DB",
				fieldSchema:               "PUBLIC",
			},
			CommandLine: []string{"snow", "connection", "test"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SNOWFLAKE_ACCOUNT":         "myorg-myaccount",
					"SNOWFLAKE_USER":            "jdoe",
					"SNOWFLAKE_AUTHENTICATOR":   "SNOWFLAKE_JWT",
					"SNOWFLAKE_PRIVATE_KEY_RAW": exampleEncryptedPrivateKey + "\n",
					"PRIVATE_KEY_PASSPHRASE":    "passphrase1234",
					"SNOWFLAKE_ROLE":            "SYSADMIN",
					"SNOWFLAKE_WAREHOUSE":       "COMPUTE_WH",
					"SNOWFLAKE_DATABASE":        "MY_DB",
					"SNOWFLAKE_SCHEMA":          "PUBLIC",
				},
				CommandLine: []string{"snow", "connection", "test", "-x"},
			},
		},
		"double dash": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Account:    "myorg-myaccount",
				fieldname.Username:   "jdoe",
				fieldname.PrivateKey: examplePrivateKey,
			},
			CommandLine: []string{"snow", "snowpark", "execute", "function", "--", "-x"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SNOWFLAKE_ACCOUNT":         "myorg-myaccount",
					"SNOWFLAKE_USER":            "jdoe",
					"SNOWFLAKE_AUTHENTICATOR":   "SNOWFLAKE_JWT",
					"SNOWFLAKE_PRIVATE_KEY_RAW": examplePrivateKey + "\n",
				},
				CommandLine: []string{"snow", "snowpark", "execute", "function", "-x", "--", "-x"},
			},
		},
	})
}

func TestNormalizePrivateKey(t *testing.T) {
	for name, input := range map[string]string{
		"line breaks":            examplePrivateKey,
		"surrounding whitespace": "\n  " + examplePrivateKey + "\n\n",
		"escaped line breaks":    "-----BEGIN PRIVATE KEY-----\\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQEXAMPLE\\n-----END PRIVATE KEY-----",
		"spaces":                 "-----BEGIN PRIVATE KEY----- MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQEXAMPLE -----END PRIVATE KEY-----",
	} {
		t.Run(name, func(t *testing.T) {
			if got := normalizePrivateKey(input); got != examplePrivateKey+"\n" {
				t.Errorf("normalizePrivateKey() = %q, want %q", got, examplePrivateKey+"\n")
			}
		})
	}
}

func TestCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, Credentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"SNOWFLAKE_ACCOUNT":         "myorg-myaccount",
				"SNOWFLAKE_USER":            "jdoe",
				"SNOWFLAKE_PRIVATE_KEY_RAW": examplePrivateKey,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:    "myorg-myaccount",
						fieldname.Username:   "jdoe",
						fieldname.PrivateKey: examplePrivateKey,
					},
				},
			},
		},
		"environment with private key file": {
			Environment: map[string]string{
				"SNOWFLAKE_ACCOUNT":          "myorg-myaccount",
				"SNOWFLAKE_USER":             "jdoe",
				"SNOWFLAKE_PRIVATE_KEY_FILE": "~/rsa_key.p8",
				"SNOWFLAKE_ROLE":             "SYSADMIN",
			},
			Files: map[string]string{
				"~/rsa_key.p8": plugintest.LoadFixture(t, "rsa_key.p8"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:    "myorg-myaccount",
						fieldname.Username:   "jdoe",
						fieldname.PrivateKey: examplePrivateKey,
						fieldRole:            "SYSADMIN",
					},
				},
			},
		},
		"environment with private key path takes precedence over raw key": {
			Environment: map[string]string{
				"SNOWFLAKE_ACCOUNT":          "myorg-myaccount",
				"SNOWFLAKE_USER":             "jdoe",
				"SNOWFLAKE_PRIVATE_KEY_PATH": "~/rsa_key.p8",
				"SNOWFLAKE_PRIVATE_KEY_RAW":  exampleEncryptedPrivateKey,
				"PRIVATE_KEY_PASSPHRASE":     "passphrase1234",
			},
			Files: map[string]string{
				"~/rsa_key.p8": plugintest.LoadFixture(t, "rsa_key.p8"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:         "myorg-myaccount",
						fieldname.Username:        "jdoe",
						fieldname.PrivateKey:      examplePrivateKey,
						fieldPrivateKeyPassphrase: "passphrase1234",
					},
				},
			},
		},
		"environment with missing private key file": {
			Environment: map[string]string{
				"SNOWFLAKE_ACCOUNT":          "myorg-myaccount",
				"SNOWFLAKE_USER":             "jdoe",
				"SNOWFLAKE_PRIVATE_KEY_FILE": "~/missing.p8",
			},
			ExpectedCandidates: []sdk.ImportCandidate{},
		},
		"config.toml with private_key_file": {
			Files: map[string]string{
				"~/.snowflake/config.toml": plugintest.LoadFixture(t, "config.toml"),
				"~/.snowflake/rsa_key.p8":  plugintest.LoadFixture(t, "rsa_key.p8"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:    "myorg-myaccount",
						fieldname.Username:   "jdoe",
						fieldname.PrivateKey: examplePrivateKey,
						fieldRole:            "SYSADMIN",
						fieldWarehouse:       "COMPUTE_WH",
					},
					NameHint: "jwt",
				},
			},
		},
		"connections.toml with private_key_raw": {
			Files: map[string]string{
				"~/.snowflake/connections.toml": plugintest.LoadFixture(t, "connections.toml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:         "myorg-myaccount",
						fieldname.Username:        "jdoe",
						fieldname.PrivateKey:      exampleEncryptedPrivateKey,
						fieldPrivateKeyPassphrase: "passphrase1234",
						fieldname.Database:        "MY_DB",
						fieldSchema:               "PUBLIC",
					},
				},
			},
		},
		"macOS config location": {
			OS: "darwin",
			Files: map[string]string{
				"~/Library/Application Support/snowflake/connections.toml": plugintest.LoadFixture(t, "connections.toml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:         "myorg-myaccount",
						fieldname.Username:        "jdoe",
						fieldname.PrivateKey:      exampleEncryptedPrivateKey,
						fieldPrivateKeyPassphrase: "passphrase1234",
						fieldname.Database:        "MY_DB",
						fieldSchema:               "PUBLIC",
					},
				},
			},
		},
	})
}
