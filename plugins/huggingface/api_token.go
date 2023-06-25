package huggingface

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIToken,
		DocsURL:       sdk.URL("https://huggingface.co/docs/huggingface_hub/quick-start"), 
		ManagementURL: sdk.URL("https://huggingface.co/settings/tokens"), 
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to HuggingFace.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 37,
					Prefix: "hf_", 
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer:           importer.TryEnvVarPair(defaultEnvVarMapping)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"HUGGING_FACE_HUB_TOKEN": fieldname.Token, 
}

