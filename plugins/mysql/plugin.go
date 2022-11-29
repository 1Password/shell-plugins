package mysql

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "mysql",
		Platform: schema.PlatformInfo{
			Name:     "MySQL",
			Homepage: sdk.URL("https://mysql.com"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			Mysql(),
		},
	}
}
