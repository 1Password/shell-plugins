package binance

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "binance",
		Platform: schema.PlatformInfo{
			Name:     "Binance",
			Homepage: sdk.URL("https://binance.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			BinanceCLI(),
		},
	}
}
