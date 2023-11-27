package binance

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func BinanceCLI() schema.Executable {
	return schema.Executable{
		Name:    "Binance CLI",
		Runs:    []string{"binance-cli"},
		DocsURL: sdk.URL("https://github.com/binance/binance-cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("t"),
			needsauth.NotWhenContainsArgs("i"),
			needsauth.NotWhenContainsArgs("book"),
			needsauth.NotWhenContainsArgs("at"),
			needsauth.NotWhenContainsArgs("k"),
			needsauth.NotWhenContainsArgs("ap"),
			needsauth.NotWhenContainsArgs("ticker"),
			needsauth.NotWhenContainsArgs("price"),
			needsauth.NotWhenContainsArgs("bt"),
			needsauth.NotWhenContainsArgs("listen"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
