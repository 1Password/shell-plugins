package contentful

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "contentful",
		Platform: schema.PlatformInfo{
			Name:     "Contentful",
			Homepage: sdk.URL("https://contentful.com"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			ContentfulCLI(),
		},
	}
}
