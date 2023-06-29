package cohere

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"fmt"
	"encoding/json"
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://docs.cohere.com/reference/config"),
		ManagementURL: sdk.URL("https://dashboard.cohere.ai/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Authtoken,
				MarkdownDescription: "JWT used to authenticate to Cohere stored in disk.",
				Secret:              true,
			},
			{
				Name:                fieldname.Email,
				MarkdownDescription: "Email used to authenticate to Cohere.",
				Optional:              true,
			},
			{
				Name:                fieldname.URL,
				MarkdownDescription: "URL of the operator server",
				Optional:              true,
			},
			
		},
		DefaultProvisioner: provision.TempFile(
			mysqlConfig, 
			provision.Filename("config"), 
			provision.AtFixedPath("~/.command/")),
		Importer: importer.TryAll(
			TryCohereConfigFile(),
		)}
}



// TODO: Check if the platform stores the Credentials in a local config file, and if so,
// implement the function below to add support for importing it.


func TryCohereConfigFile() sdk.Importer {
	return importer.TryFile("~/command/config", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.Contexts[config.CurrentURL].JWT == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.URL: config.CurrentURL,
				fieldname.Authtoken: config.Contexts[config.CurrentURL].JWT,
				fieldname.Email: config.Contexts[config.CurrentURL].Email,
				
			},
		})
	})
}
func mysqlConfig(in sdk.ProvisionInput) ([]byte, error) {
var currentURL, jwt, email string


if value, ok := in.ItemFields[fieldname.URL]; ok {
	currentURL = value
}

if value, ok := in.ItemFields[fieldname.Authtoken]; ok {
	jwt = value
}

if value, ok := in.ItemFields[fieldname.Email]; ok {
	email = value
}

data := Config{
	CurrentURL: currentURL,
	Contexts: map[string]APISettings{
		currentURL: {
			JWT:   jwt,
			Email: email,
		},
	},
	
}

jsonData, err := json.MarshalIndent(data, "", "\t")
if err != nil {
	fmt.Println("Error:", err)
	
}
return []byte(jsonData), nil
}

// TODO: Implement the config file schema
type Config struct {
	CurrentURL string                 `json:"CurrentURL"`
	Contexts   map[string]APISettings `json:"Contexts"`
}

type APISettings struct {
	JWT   string `json:"JWT"`
	Email string `json:"Email"`
}

