package civo

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/mitchellh/go-homedir"
)

func APIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		DocsURL:       sdk.URL("https://www.civo.com/docs/account/api-keys"),
		ManagementURL: sdk.URL("https://dashboard.civo.com/security"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Civo.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 50,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.APIKeyID,
				MarkdownDescription: "API Name to identify the API Key.",
			},
			{
				Name:                fieldname.DefaultRegion,
				MarkdownDescription: "The default region to use for this API Key.",
				Optional:            true,
			},
		},
		// DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		DefaultProvisioner: provision.TempFile(configFile,
			provision.Filename(".civo.json"),
			provision.AddArgs(
				"--config", "{{.Path}}",
			),
			provision.AtFixedPath("~/.civo.json"),
		),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryCivoConfigFile(),
		)}
}
func configFile(in sdk.ProvisionInput) ([]byte, error) {
	apiKey := in.ItemFields[fieldname.APIKey]
	apiKeyID := in.ItemFields[fieldname.APIKeyID]
	defaultRegion := in.ItemFields[fieldname.DefaultRegion]

	// Check if the file already exists
	filePath := "~/.civo.json"
	exists, err := fileExists(filePath)
	if err != nil {
		return nil, err
	}

	if exists {
		// Read the existing file
		contents, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var config Config
		if err := json.Unmarshal(contents, &config); err != nil {
			return nil, err
		}

		// Update the config with the new values
		config.Properties[apiKeyID] = json.RawMessage(`"` + apiKey + `"`)
		config.Meta.CurrentAPIKey = apiKeyID
		config.Meta.DefaultRegion = defaultRegion

		return json.MarshalIndent(config, "", "  ")
	}

	// Create a new config
	config := Config{
		Properties: map[string]json.RawMessage{
			apiKeyID: json.RawMessage(`"` + apiKey + `"`),
		},
		Meta: struct {
			Admin               bool   `json:"admin"`
			CurrentAPIKey       string `json:"current_apikey"`
			DefaultRegion       string `json:"default_region"`
			LatestReleaseCheck  string `json:"latest_release_check"`
			URL                 string `json:"url"`
			LastCommandExecuted string `json:"last_command_executed"`
		}{
			CurrentAPIKey: apiKeyID,
			DefaultRegion: defaultRegion,
		},
	}

	return json.MarshalIndent(config, "", "  ")
}

func fileExists(filePath string) (bool, error) {
	expandedPath, err := homedir.Expand(filePath)
	if err != nil {
		return false, err
	}

	info, err := os.Stat(expandedPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return !info.IsDir(), nil
}

// func configFile(in sdk.ProvisionInput) ([]byte, error) {
// 	apiKey := in.ItemFields[fieldname.APIKey]
// 	apiKeyID := in.ItemFields[fieldname.APIKeyID]
// 	defaultRegion := in.ItemFields[fieldname.DefaultRegion]

// 	config := Config{
// 		Properties: map[string]json.RawMessage{
// 			apiKeyID: json.RawMessage(`"` + apiKey + `"`),
// 		},
// 		Meta: struct {
// 			Admin               bool   `json:"admin"`
// 			CurrentAPIKey       string `json:"current_apikey"`
// 			DefaultRegion       string `json:"default_region"`
// 			LatestReleaseCheck  string `json:"latest_release_check"`
// 			URL                 string `json:"url"`
// 			LastCommandExecuted string `json:"last_command_executed"`
// 		}{
// 			CurrentAPIKey: apiKeyID,
// 			DefaultRegion: defaultRegion,
// 		},
// 	}

// 	return json.MarshalIndent(config, "", "  ")
// }

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"CIVO_API_KEY_NAME": fieldname.APIKeyID,
	"CIVO_API_KEY":      fieldname.APIKey,
}

func TryCivoConfigFile() sdk.Importer {

	return importer.TryFile("~/.civo.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return

		}

		if len(config.Properties) == 0 && config.Meta.CurrentAPIKey == "" {
			return
		}

		var apiKey string
		for key, value := range config.Properties {
			if key == config.Meta.CurrentAPIKey {
				err := json.Unmarshal(value, &apiKey)
				if err != nil {
					out.AddError(err)
					return
				}
			}
			break
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.APIKey:        apiKey,
				fieldname.APIKeyID:      config.Meta.CurrentAPIKey,
				fieldname.DefaultRegion: config.Meta.DefaultRegion,
			},
		})

	})
}

type Config struct {
	Properties map[string]json.RawMessage `json:"apikeys"`

	Meta struct {
		Admin               bool   `json:"admin"`
		CurrentAPIKey       string `json:"current_apikey"`
		DefaultRegion       string `json:"default_region"`
		LatestReleaseCheck  string `json:"latest_release_check"`
		URL                 string `json:"url"`
		LastCommandExecuted string `json:"last_command_executed"`
	} `json:"meta"`
}
