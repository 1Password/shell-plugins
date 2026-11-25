package qdrantcloud

import (
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func QdrantCloudCLI() schema.Executable {
	return schema.Executable{
		Name:    "Qdrant Cloud CLI",
		Runs:    []string{"qcloud"},
		DocsURL: sdk.URL("https://qdrant.tech/documentation/cloud-cli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForCommand("version"),
			needsauth.NotForCommand("context"),
			needsauth.NotForCommand("completion"),
			needsauth.NotForCommand("self-upgrade"),
			needsauth.NotForCommand("__complete"),
			needsauth.NotForCommand("__completeNoDesc"),
			notWithContextFlag(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.ManagementKey,
			},
		},
	}
}

// notWithContextFlag opts out of authentication when a context is picked explicitly with --context,
// since provisioned env vars would otherwise take precedence over that context's settings.
func notWithContextFlag() sdk.NeedsAuthentication {
	return func(in sdk.NeedsAuthenticationInput) bool {
		for _, arg := range in.CommandArgs {
			if arg == "--context" || strings.HasPrefix(arg, "--context=") {
				return false
			}
		}
		return true
	}
}
