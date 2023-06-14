package redis

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func RedisCLI() schema.Executable {
	return schema.Executable{
		Name:    "Redis CLI",
		Runs:    []string{"redis-cli"},
		DocsURL: sdk.URL("https://redis.io/docs/ui/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotWhenContainsArgs("-u"),
			needsauth.NotWhenContainsArgs("--user"), // skip when a custom username is provided
			needsauth.NotWhenContainsArgs("-a"),     // skip when a custom password is provided
			needsauth.NotWhenContainsArgs("-h"),     // skip when a custom host is provided
			needsauth.NotWhenContainsArgs("-p"),     // skip when a custom port is provided
			needsauth.NotWhenContainsArgs("--help"),
			needsauth.NotForVersion(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.UserCredentials,
				Provisioner: redisProvisioner(),
			},
		},
	}
}
