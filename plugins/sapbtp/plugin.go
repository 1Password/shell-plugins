package sapbtp

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "sapbtp",
		Platform: schema.PlatformInfo{
			Name:     "SAP BTP",
			Homepage: sdk.URL("https://help.sap.com/docs/BTP"),
		},
		Credentials: []schema.CredentialType{
			Credentials(),
		},
		Executables: []schema.Executable{
			SAPBTPCLI(),
		},
	}
}
