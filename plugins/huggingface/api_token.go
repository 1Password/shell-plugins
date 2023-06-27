package huggingface

import (
	"context"

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
			{
				Name:                fieldname.OrgURL,
				MarkdownDescription: "Endpoint used to connect to HuggingFace CLI",
				Optional:            true,
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.Website,
				MarkdownDescription: "HF Inference Endpoint used to connect to HuggingFace CLI",
				Optional:            true,
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryHFtokenFile(),
		),}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"HUGGING_FACE_HUB_TOKEN": fieldname.Token, 
	"HF_ENDPOINT": fieldname.OrgURL,
	"HF_INFERENCE_ENDPOINT": fieldname.Endpoint,

}


func TryHFtokenFile() sdk.Importer {
	return importer.TryFile("~/.cache/huggingface/token", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		fileData := string(contents)

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Token: fileData ,
			},
		})
	})

}
