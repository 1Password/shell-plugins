package cockroachdb

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestDatabaseCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, DatabaseCredentials().Importer, map[string]plugintest.ImportCase{
		"environment variables - complete": {
			Environment: map[string]string{
				"COCKROACH_HOST":     "localhost",
				"COCKROACH_PORT":     "26257",
				"COCKROACH_USER":     "root",
				"COCKROACH_PASSWORD": "password123",
				"COCKROACH_DATABASE": "defaultdb",
				"COCKROACH_INSECURE": "1",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Host:     "localhost",
						fieldname.Port:     "26257",
						fieldname.User:     "root",
						fieldname.Password: "password123",
						fieldname.Database: "defaultdb",
						"insecure":         "1",
					},
				},
			},
		},
		"environment variables - minimal": {
			Environment: map[string]string{
				"COCKROACH_HOST": "cockroach.example.com",
				"COCKROACH_USER": "admin",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Host: "cockroach.example.com",
						fieldname.User: "admin",
					},
				},
			},
		},
		"environment variables - production secure": {
			Environment: map[string]string{
				"COCKROACH_HOST":     "prod-cluster.cockroachlabs.cloud",
				"COCKROACH_PORT":     "26257",
				"COCKROACH_USER":     "produser",
				"COCKROACH_PASSWORD": "securepass123",
				"COCKROACH_DATABASE": "proddb",
				"COCKROACH_INSECURE": "0",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Host:     "prod-cluster.cockroachlabs.cloud",
						fieldname.Port:     "26257",
						fieldname.User:     "produser",
						fieldname.Password: "securepass123",
						fieldname.Database: "proddb",
						"insecure":         "0",
					},
				},
			},
		},
	})
}

func TestDatabaseCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, DatabaseCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"local development - insecure": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Host:     "localhost",
				fieldname.Port:     "26257",
				fieldname.User:     "root",
				fieldname.Password: "password123",
				fieldname.Database: "defaultdb",
				"insecure":         "1",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"COCKROACH_HOST":     "localhost",
					"COCKROACH_PORT":     "26257",
					"COCKROACH_USER":     "root",
					"COCKROACH_PASSWORD": "password123",
					"COCKROACH_DATABASE": "defaultdb",
					"COCKROACH_INSECURE": "1",
				},
			},
		},
		"production - secure": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Host:     "prod-cluster.cockroachlabs.cloud",
				fieldname.Port:     "26257",
				fieldname.User:     "produser",
				fieldname.Password: "securepass123",
				fieldname.Database: "proddb",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"COCKROACH_HOST":     "prod-cluster.cockroachlabs.cloud",
					"COCKROACH_PORT":     "26257",
					"COCKROACH_USER":     "produser",
					"COCKROACH_PASSWORD": "securepass123",
					"COCKROACH_DATABASE": "proddb",
				},
			},
		},
		"minimal configuration": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Host: "cockroach.example.com",
				fieldname.User: "admin",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"COCKROACH_HOST": "cockroach.example.com",
					"COCKROACH_USER": "admin",
				},
			},
		},
	})
}

// TestCockroachSQLExecutable tests the cockroach sql executable configuration
func TestCockroachSQLExecutable(t *testing.T) {
	plugin := New()

	// Find the cockroach sql executable
	var cockroachSQL *schema.Executable
	for _, exec := range plugin.Executables {
		if exec.Name == "cockroach" {
			cockroachSQL = &exec
			break
		}
	}

	if cockroachSQL == nil {
		t.Fatal("cockroach sql executable not found in plugin")
	}

	// Test that it uses database credentials
	if len(cockroachSQL.Uses) != 1 {
		t.Errorf("Expected 1 credential usage, got %d", len(cockroachSQL.Uses))
	}

	if cockroachSQL.Uses[0].Name != credname.DatabaseCredentials {
		t.Errorf("Expected DatabaseCredentials, got %s", cockroachSQL.Uses[0].Name)
	}
}

// TestPluginValidation tests that the plugin passes all validation checks
func TestPluginValidation(t *testing.T) {
	plugin := New()

	// Basic plugin validation
	if plugin.Name != "cockroachdb" {
		t.Errorf("Expected plugin name 'cockroachdb', got '%s'", plugin.Name)
	}

	if len(plugin.Credentials) != 1 {
		t.Errorf("Expected 1 credential type, got %d", len(plugin.Credentials))
	}

	if len(plugin.Executables) != 1 {
		t.Errorf("Expected 1 executable, got %d", len(plugin.Executables))
	}
}
