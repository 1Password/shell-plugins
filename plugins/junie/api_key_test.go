package junie

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"junie api key": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey: "perm-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"JUNIE_API_KEY": "perm-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
				},
			},
		},
		"does not provision provider api keys": {
			ItemFields: map[sdk.FieldName]string{
				sdk.FieldName("Anthropic API Key"):  "sk-ant-api03-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
				sdk.FieldName("OpenAI API Key"):     "sk-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
				sdk.FieldName("Google API Key"):     "AI-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
				sdk.FieldName("Grok API Key"):       "xai-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
				sdk.FieldName("OpenRouter API Key"): "sk-or-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{},
			},
		},
		"does not provision project or task runtime inputs": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Project:     "/tmp/example-project",
				sdk.FieldName("Task"): "Review and fix any code quality issues in the latest commit",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"junie api key environment": {
			Environment: map[string]string{
				"JUNIE_API_KEY": "perm-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "perm-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
					},
				},
			},
		},
		"provider api key environment": {
			Environment: map[string]string{
				"JUNIE_ANTHROPIC_API_KEY":  "sk-ant-api03-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
				"JUNIE_OPENAI_API_KEY":     "sk-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
				"JUNIE_GOOGLE_API_KEY":     "AI-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
				"JUNIE_GROK_API_KEY":       "xai-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
				"JUNIE_OPENROUTER_API_KEY": "sk-or-yEyY18xzH5IiiORdCDzstp1h2xrxCydfh9tjFveUyEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{},
		},
	})
}
