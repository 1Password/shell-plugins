package curl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "curl",
		Platform: schema.PlatformInfo{
			Name:     "cURL",
			Homepage: sdk.URL("https://curl.se/"),
		},
		Executables: []schema.Executable{
			CurlCLI(),
		},
	}
}
