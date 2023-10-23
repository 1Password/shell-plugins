package huggingface

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func HuggingFaceCLI() schema.Executable {
	return schema.Executable{
		Name:    "HuggingFace CLI",
		Runs:    []string{"huggingface-cli"},
		DocsURL: sdk.URL("https://huggingface.co/docs/huggingface_hub/quick-start"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("login"),
			needsauth.NotWhenContainsArgs("logout"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
