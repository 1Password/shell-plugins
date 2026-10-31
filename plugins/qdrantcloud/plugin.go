package qdrantcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "qdrantcloud",
		Platform: schema.PlatformInfo{
			Name:     "Qdrant Cloud",
			Homepage: sdk.URL("https://qdrant.tech"),
		},
		Credentials: []schema.CredentialType{
			ManagementKey(),
		},
		Executables: []schema.Executable{
			QdrantCloudCLI(),
		},
	}
}
