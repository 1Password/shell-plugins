package redis

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "redis",
		Platform: schema.PlatformInfo{
			Name:     "Redis",
			Homepage: sdk.URL("https://redis.io/"),
		},
		Credentials: []schema.CredentialType{
			UserCredentials(),
		},
		Executables: []schema.Executable{
			RedisCLI(),
		},
	}
}
