package llm

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				OpenAIFieldName:    "sk-proj-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
				AnthropicFieldName: "sk-ant-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"OPENAI_API_KEY":    "sk-proj-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
					"ANTHROPIC_API_KEY": "sk-ant-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"OPENAI_API_KEY":    "sk-proj-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
				"ANTHROPIC_API_KEY": "sk-ant-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						OpenAIFieldName:    "sk-proj-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
						AnthropicFieldName: "sk-ant-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
					},
				},
			},
		},
		"config file mac": {
			OS: "darwin",
			Files: map[string]string{
				"~/Library/Application Support/io.datasette.llm/keys.json": plugintest.LoadFixture(t, "keys.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						OpenAIFieldName:    "sk-proj-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
						AnthropicFieldName: "sk-ant-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
					},
				},
			},
		},
		"config file linux": {
			OS: "linux",
			Files: map[string]string{
				"~/.config/io.datasette.llm/keys.json": plugintest.LoadFixture(t, "keys.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						OpenAIFieldName:    "sk-proj-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
						AnthropicFieldName: "sk-ant-ysT1SpYOenNu805nCf3yUYIbNAfvHSNzR0rx2WGRHEXAMPLE",
					},
				},
			},
		},
	})
}
