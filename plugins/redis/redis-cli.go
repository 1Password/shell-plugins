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
		DocsURL: sdk.URL("https://redis.io/docs/manual/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("-u"),
			needsauth.NotWhenContainsArgs("-a"),
			needsauth.NotWhenContainsArgs("-v"),
			needsauth.NotWhenContainsArgs("--version"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Password,
			},
		},
	}
}
