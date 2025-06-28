package scorecard

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func generateTestPrivateKey(t *testing.T) string {
	// Generate a test RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test private key: %v", err)
	}

	// Encode to PEM format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return string(privateKeyPEM)
}

func TestSecretKeyProvisioner(t *testing.T) {
	testPrivateKey := generateTestPrivateKey(t)

	plugintest.TestProvisioner(t, SecretKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Key: testPrivateKey,
				"App ID":         "123456",
				"Installation ID": "7890123",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GITHUB_APP_ID":              "123456",
					"GITHUB_APP_INSTALLATION_ID": "7890123",
					"GITHUB_APP_KEY_PATH":        "/tmp/github-app-key.pem",
				},
				Files: map[string]sdk.OutputFile{
					"/tmp/github-app-key.pem": {Contents: []byte(testPrivateKey)},
				},
			},
		},
		"partial fields": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Key: testPrivateKey,
				"App ID": "123456",
				// Installation ID is missing to test handling of partial fields
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GITHUB_APP_ID":       "123456",
					"GITHUB_APP_KEY_PATH": "/tmp/github-app-key.pem",
					// GITHUB_APP_INSTALLATION_ID should not be set when field is missing
				},
				Files: map[string]sdk.OutputFile{
					"/tmp/github-app-key.pem": {Contents: []byte(testPrivateKey)},
				},
			},
		},
	})
}

func TestSecretKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, SecretKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"GITHUB_APP_KEY_PATH":         "/path/to/key.pem",
				"GITHUB_APP_ID":               "123456", 
				"GITHUB_APP_INSTALLATION_ID":  "7890123",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Key:     "/path/to/key.pem",
						"App ID":          "123456",
						"Installation ID": "7890123",
					},
				},
			},
		},
	})
}