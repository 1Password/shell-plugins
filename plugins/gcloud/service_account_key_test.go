package gcloud

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestServiceAccountKeyProvisioner(t *testing.T) {
	saKeyJSON := plugintest.LoadFixture(t, "service_account_key.json")
	adcJSON := plugintest.LoadFixture(t, "application_default_credentials.json")

	plugintest.TestProvisioner(t, ServiceAccountKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"service account key": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Credential: saKeyJSON,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GOOGLE_APPLICATION_CREDENTIALS": "/tmp/gcloud-credentials.json",
				},
				Files: map[string]sdk.OutputFile{
					"/tmp/gcloud-credentials.json": {Contents: []byte(saKeyJSON)},
				},
			},
		},
		"service account key with project": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Credential: saKeyJSON,
				fieldname.ProjectID:  "my-gcp-project",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GOOGLE_APPLICATION_CREDENTIALS": "/tmp/gcloud-credentials.json",
					"CLOUDSDK_CORE_PROJECT":          "my-gcp-project",
				},
				Files: map[string]sdk.OutputFile{
					"/tmp/gcloud-credentials.json": {Contents: []byte(saKeyJSON)},
				},
			},
		},
		"authorized user credential": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Credential: adcJSON,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GOOGLE_APPLICATION_CREDENTIALS": "/tmp/gcloud-credentials.json",
				},
				Files: map[string]sdk.OutputFile{
					"/tmp/gcloud-credentials.json": {Contents: []byte(adcJSON)},
				},
			},
		},
	})
}

func TestServiceAccountKeyImporter(t *testing.T) {
	saKeyJSON := plugintest.LoadFixture(t, "service_account_key.json")
	adcJSON := plugintest.LoadFixture(t, "application_default_credentials.json")

	plugintest.TestImporter(t, ServiceAccountKey().Importer, map[string]plugintest.ImportCase{
		"ADC file with service account key": {
			Files: map[string]string{
				"~/.config/gcloud/application_default_credentials.json": saKeyJSON,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Credential: saKeyJSON,
						fieldname.ProjectID:  "my-gcp-project",
					},
				},
			},
		},
		"ADC file with authorized user credential": {
			Files: map[string]string{
				"~/.config/gcloud/application_default_credentials.json": adcJSON,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Credential: adcJSON,
					},
				},
			},
		},
	})
}

func TestGCloudCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, GCloudCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no args": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
		"help flag": {
			Args:              []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"version flag": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"auth login": {
			Args:              []string{"auth", "login"},
			ExpectedNeedsAuth: false,
		},
		"auth list": {
			Args:              []string{"auth", "list"},
			ExpectedNeedsAuth: false,
		},
		"config set": {
			Args:              []string{"config", "set", "project", "my-project"},
			ExpectedNeedsAuth: false,
		},
		"info": {
			Args:              []string{"info"},
			ExpectedNeedsAuth: false,
		},
		"components list": {
			Args:              []string{"components", "list"},
			ExpectedNeedsAuth: false,
		},
		"compute instances list": {
			Args:              []string{"compute", "instances", "list"},
			ExpectedNeedsAuth: true,
		},
		"storage ls": {
			Args:              []string{"storage", "ls"},
			ExpectedNeedsAuth: true,
		},
		"projects list": {
			Args:              []string{"projects", "list"},
			ExpectedNeedsAuth: true,
		},
	})
}

func TestGsutilCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, GsutilCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no args": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
		"help flag": {
			Args:              []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"version flag": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"ls": {
			Args:              []string{"ls"},
			ExpectedNeedsAuth: true,
		},
	})
}

func TestBqCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, BqCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no args": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
		"help flag": {
			Args:              []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"version flag": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"query": {
			Args:              []string{"query", "SELECT 1"},
			ExpectedNeedsAuth: true,
		},
	})
}
