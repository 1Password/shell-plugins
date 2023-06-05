package redis

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func RedisCLI() schema.Executable {
	return schema.Executable{
		Name:    "Redis CLI",
		Runs:    []string{"redis-cli"},
		DocsURL: sdk.URL("https://redis.io/docs/ui/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotWhenContainsArgs("-u"),
			needsauth.NotWhenContainsArgs("--user"),
			needsauth.NotWhenContainsArgs("-a"),
			needsauth.NotWhenContainsArgs("--help"),
			needsauth.NotForVersion(),
		),
		Uses: []schema.CredentialUsage{
			{
				SelectFrom: &schema.CredentialSelection{
					IncludeAllCredentials: false,
					AllowMultiple:         false,
				},
				Optional:    false,
				Provisioner: EnvVarFlags(flagsToProvision),
			},
		},
	}
}
