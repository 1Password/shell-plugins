package sentry

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "sentry-cli",
		Platform: schema.PlatformInfo{
			Name:     "Sentry",
			Homepage: sdk.URL("https://sentry.io"),
		},
		Credentials: []schema.CredentialType{
			AuthToken(),
		},
		Executables: []schema.Executable{
			Executable_sentry_cli(),
		},
	}
}
