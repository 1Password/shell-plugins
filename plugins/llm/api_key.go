package llm

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

// LLM enables you to query many APIs with a single command.
// Therefore, we have an optional field for each API key in the 1Password entry.

const (
	AnthropicFieldName  = sdk.FieldName("Anthropic")
	AnyscaleFieldName   = sdk.FieldName("Anyscale")
	CohereFieldName     = sdk.FieldName("Cohere")
	FireworksFieldName  = sdk.FieldName("Fireworks")
	GeminiFieldName     = sdk.FieldName("Gemini")
	GroqFieldName       = sdk.FieldName("Groq")
	MistralFieldName    = sdk.FieldName("Mistral")
	OpenAIFieldName     = sdk.FieldName("OpenAI")
	OpenRouterFieldName = sdk.FieldName("OpenRouter")
	PALMFieldName       = sdk.FieldName("PaLM")
	PerplexityFieldName = sdk.FieldName("Perplexity")
	RekaFieldName       = sdk.FieldName("Reka")
	ReplicateFieldName  = sdk.FieldName("Replicate")
	TogetherFieldName   = sdk.FieldName("Together")
)

var defaultValueComposition = &schema.ValueComposition{
	Charset: schema.Charset{
		Uppercase: true,
		Lowercase: true,
		Digits:    true,
		Symbols:   true,
	},
}

var availableServices = []ServiceDefinition{
	{
		FieldName:           AnthropicFieldName,
		EnvVarName:          "ANTHROPIC_API_KEY",
		MarkdownDescription: "API Key for Anthropic.",
		Composition: &schema.ValueComposition{
			Prefix: "sk-ant-",
			Charset: schema.Charset{
				Uppercase: true,
				Lowercase: true,
				Digits:    true,
				Symbols:   true,
			},
		},
		ConfigFieldFunc: func(c Config) string { return c.Claude },
	},
	{
		FieldName:           AnyscaleFieldName,
		EnvVarName:          "LLM_ANYSCALE_ENDPOINTS_KEY",
		MarkdownDescription: "API Key for Anyscale.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.Anyscale },
	},
	{
		FieldName:           CohereFieldName,
		EnvVarName:          "COHERE_API_KEY",
		MarkdownDescription: "API Key for Cohere.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.Cohere },
	},
	{
		FieldName:           FireworksFieldName,
		EnvVarName:          "LLM_FIREWORKS_KEY",
		MarkdownDescription: "API Key for Fireworks.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.Fireworks },
	},
	{
		FieldName:           GeminiFieldName,
		EnvVarName:          "LLM_GEMINI_KEY",
		MarkdownDescription: "API Key for Google’s Gemini.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.Gemini },
	},
	{
		FieldName:           GroqFieldName,
		EnvVarName:          "LLM_GROQ_KEY",
		MarkdownDescription: "API Key for Groq.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.Groq },
	},
	{
		FieldName:           MistralFieldName,
		EnvVarName:          "LLM_MISTRAL_KEY",
		MarkdownDescription: "API Key for Mistral.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.Mistral },
	},
	{
		FieldName:           OpenAIFieldName,
		EnvVarName:          "OPENAI_API_KEY",
		MarkdownDescription: "API Key for OpenAI.",
		Composition: &schema.ValueComposition{
			Prefix: "sk-",
			Charset: schema.Charset{
				Uppercase: true,
				Lowercase: true,
				Digits:    true,
			},
		},
		ConfigFieldFunc: func(c Config) string { return c.OpenAI },
	},
	{
		FieldName:           OpenRouterFieldName,
		EnvVarName:          "LLM_OPENROUTER_KEY",
		MarkdownDescription: "API Key for OpenRouter.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.OpenRouter },
	},
	{
		FieldName:           PALMFieldName,
		EnvVarName:          "PALM_API_KEY",
		MarkdownDescription: "API Key for PALM.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.PALM },
	},
	{
		FieldName:           PerplexityFieldName,
		EnvVarName:          "PERPLEXITY_API_KEY",
		MarkdownDescription: "API Key for Perplexity.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.Perplexity },
	},
	{
		FieldName:           RekaFieldName,
		EnvVarName:          "LLM_REKA_KEY",
		MarkdownDescription: "API Key for Reka.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.Reka },
	},
	{
		FieldName:           ReplicateFieldName,
		EnvVarName:          "REPLICATE_API_KEY",
		MarkdownDescription: "API Key for Replicate.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.Replicate },
	},
	{
		FieldName:           TogetherFieldName,
		EnvVarName:          "TOGETHER_API_KEY",
		MarkdownDescription: "API Key for Together.",
		Composition:         defaultValueComposition,
		ConfigFieldFunc:     func(c Config) string { return c.Together },
	},
}

func APIKey() schema.CredentialType {
	// Create env variable mapping for each LLM service
	var defaultEnvVarMapping = make(map[string]sdk.FieldName)
	for _, field := range availableServices {
		if field.EnvVarName != "" {
			defaultEnvVarMapping[field.EnvVarName] = field.FieldName
		}
	}

	// Create schema fields for each LLM service
	var schemaFields []schema.CredentialField
	for _, field := range availableServices {
		schemaFields = append(schemaFields, schema.CredentialField{
			Name:                field.FieldName,
			MarkdownDescription: field.MarkdownDescription,
			Secret:              true,
			Optional:            true,
			Composition:         field.Composition,
		})
	}

	return schema.CredentialType{
		Name:               credname.APIKey,
		DocsURL:            sdk.URL("https://llm.datasette.io/en/stable/setup.html"),
		Fields:             schemaFields,
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			importer.MacOnly(TryLLMConfigFile("~/Library/Application Support/io.datasette.llm/keys.json")),
			importer.LinuxOnly(TryLLMConfigFile("~/.config/io.datasette.llm/keys.json")),
		)}
}

func TryLLMConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		// Add candidates for each service that has a value in the config file
		candidateFields := make(map[sdk.FieldName]string)
		for _, field := range availableServices {
			var configValue string = field.ConfigFieldFunc(config)
			if configValue != "" {
				candidateFields[field.FieldName] = configValue
			}
		}

		if len(candidateFields) == 0 {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: candidateFields,
		})
	})
}

type ServiceDefinition struct {
	FieldName           sdk.FieldName
	ConfigFileFieldName string
	EnvVarName          string
	MarkdownDescription string
	Composition         *schema.ValueComposition
	ConfigFieldFunc     func(Config) string
}

type Config struct {
	Anyscale   string `json:"anyscale-endpoints"`
	Claude     string `json:"claude"`
	Cohere     string `json:"cohere"`
	Fireworks  string `json:"fireworks"`
	Gemini     string `json:"gemini"`
	Groq       string `json:"groq"`
	Mistral    string `json:"mistral"`
	OpenAI     string `json:"openai"`
	OpenRouter string `json:"openrouter"`
	PALM       string `json:"palm"`
	Perplexity string `json:"perplexity"`
	Reka       string `json:"reka"`
	Replicate  string `json:"replicate"`
	Together   string `json:"together"`
}
