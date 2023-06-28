package huggingface

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "huggingface",
		Platform: schema.PlatformInfo{
			Name:     "Hugging Face",
			Homepage: sdk.URL("https://huggingface.co"),
		},
		Credentials: []schema.CredentialType{
			User_Access_Token(),
		},
		Executables: []schema.Executable{
			HuggingFaceCLI(),
		},
	}
}
