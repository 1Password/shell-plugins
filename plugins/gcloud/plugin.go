package gcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "gcloud",
		Platform: schema.PlatformInfo{
			Name:     "Google Cloud Platform",
			Homepage: sdk.URL("https://cloud.google.com"),
		},
		Credentials: []schema.CredentialType{
			ServiceAccountKey(),
		},
		Executables: []schema.Executable{
			GCloudCLI(),
			GsutilCLI(),
			BqCLI(),
		},
	}
}
