package qdrantcloud

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestManagementKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, ManagementKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.ManagementKey: "9f4c2a7e-1b3d-4e5f-8a6b-7c8d9e0f1a2b|qX7vR2kP9mT4wZ8nB3cJ6hL1sD5fGa8Yt2Wk4Ns6Pd0Rm3VhEXAMPLE",
				fieldname.AccountID:     "3e8b1c4d-5f6a-4b7c-9d0e-1f2a3b4c5d6e",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"QDRANT_CLOUD_API_KEY":    "9f4c2a7e-1b3d-4e5f-8a6b-7c8d9e0f1a2b|qX7vR2kP9mT4wZ8nB3cJ6hL1sD5fGa8Yt2Wk4Ns6Pd0Rm3VhEXAMPLE",
					"QDRANT_CLOUD_ACCOUNT_ID": "3e8b1c4d-5f6a-4b7c-9d0e-1f2a3b4c5d6e",
				},
			},
		},
	})
}

func TestManagementKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, ManagementKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"QDRANT_CLOUD_API_KEY":    "9f4c2a7e-1b3d-4e5f-8a6b-7c8d9e0f1a2b|qX7vR2kP9mT4wZ8nB3cJ6hL1sD5fGa8Yt2Wk4Ns6Pd0Rm3VhEXAMPLE",
				"QDRANT_CLOUD_ACCOUNT_ID": "3e8b1c4d-5f6a-4b7c-9d0e-1f2a3b4c5d6e",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.ManagementKey: "9f4c2a7e-1b3d-4e5f-8a6b-7c8d9e0f1a2b|qX7vR2kP9mT4wZ8nB3cJ6hL1sD5fGa8Yt2Wk4Ns6Pd0Rm3VhEXAMPLE",
						fieldname.AccountID:     "3e8b1c4d-5f6a-4b7c-9d0e-1f2a3b4c5d6e",
					},
				},
			},
		},
		"environment (key only)": {
			Environment: map[string]string{
				"QDRANT_CLOUD_API_KEY": "9f4c2a7e-1b3d-4e5f-8a6b-7c8d9e0f1a2b|qX7vR2kP9mT4wZ8nB3cJ6hL1sD5fGa8Yt2Wk4Ns6Pd0Rm3VhEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.ManagementKey: "9f4c2a7e-1b3d-4e5f-8a6b-7c8d9e0f1a2b|qX7vR2kP9mT4wZ8nB3cJ6hL1sD5fGa8Yt2Wk4Ns6Pd0Rm3VhEXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.config/qcloud/config.yaml": plugintest.LoadFixture(t, "config.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.ManagementKey: "9f4c2a7e-1b3d-4e5f-8a6b-7c8d9e0f1a2b|qX7vR2kP9mT4wZ8nB3cJ6hL1sD5fGa8Yt2Wk4Ns6Pd0Rm3VhEXAMPLE",
						fieldname.AccountID:     "3e8b1c4d-5f6a-4b7c-9d0e-1f2a3b4c5d6e",
					},
					NameHint: "production",
				},
			},
		},
		"config file (yml)": {
			Files: map[string]string{
				"~/.config/qcloud/config.yml": plugintest.LoadFixture(t, "config.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.ManagementKey: "9f4c2a7e-1b3d-4e5f-8a6b-7c8d9e0f1a2b|qX7vR2kP9mT4wZ8nB3cJ6hL1sD5fGa8Yt2Wk4Ns6Pd0Rm3VhEXAMPLE",
						fieldname.AccountID:     "3e8b1c4d-5f6a-4b7c-9d0e-1f2a3b4c5d6e",
					},
					NameHint: "production",
				},
			},
		},
		"config file (QDRANT_CLOUD_CONFIG)": {
			Environment: map[string]string{
				"QDRANT_CLOUD_CONFIG": "/work/qcloud.yaml",
			},
			Files: map[string]string{
				"/work/qcloud.yaml": plugintest.LoadFixture(t, "config.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.ManagementKey: "9f4c2a7e-1b3d-4e5f-8a6b-7c8d9e0f1a2b|qX7vR2kP9mT4wZ8nB3cJ6hL1sD5fGa8Yt2Wk4Ns6Pd0Rm3VhEXAMPLE",
						fieldname.AccountID:     "3e8b1c4d-5f6a-4b7c-9d0e-1f2a3b4c5d6e",
					},
					NameHint: "production",
				},
			},
		},
		"config file (api_key_command only)": {
			Files: map[string]string{
				"~/.config/qcloud/config.yaml": "contexts:\n  - name: default\n    api_key_command: op read op://Private/qdrant/credential\n",
			},
			ExpectedCandidates: nil,
		},
	})
}

func TestQdrantCloudCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, QdrantCloudCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no for no args":              {Args: []string{}, ExpectedNeedsAuth: false},
		"no for --help":               {Args: []string{"cluster", "--help"}, ExpectedNeedsAuth: false},
		"no for version":              {Args: []string{"version"}, ExpectedNeedsAuth: false},
		"no for version --json":       {Args: []string{"version", "--json"}, ExpectedNeedsAuth: false},
		"no for --context":            {Args: []string{"--context", "staging", "cluster", "list"}, ExpectedNeedsAuth: false},
		"no for --context=":           {Args: []string{"cluster", "list", "--context=staging"}, ExpectedNeedsAuth: false},
		"no for context set":          {Args: []string{"context", "set", "default"}, ExpectedNeedsAuth: false},
		"no for completion":           {Args: []string{"completion", "zsh"}, ExpectedNeedsAuth: false},
		"no for self-upgrade":         {Args: []string{"self-upgrade"}, ExpectedNeedsAuth: false},
		"no for shell completions":    {Args: []string{"__complete", "cluster", "describe", ""}, ExpectedNeedsAuth: false},
		"no for completions w/o desc": {Args: []string{"__completeNoDesc", "cluster", "describe", ""}, ExpectedNeedsAuth: false},
		"yes for account list":        {Args: []string{"account", "list"}, ExpectedNeedsAuth: true},
		"yes for cluster list":        {Args: []string{"cluster", "list"}, ExpectedNeedsAuth: true},
		"yes for cloud-provider":      {Args: []string{"cloud-provider", "list", "--json"}, ExpectedNeedsAuth: true},
	})
}
