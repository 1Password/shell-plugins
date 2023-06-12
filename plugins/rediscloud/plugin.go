package rediscloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "rediscloud",
		Platform: schema.PlatformInfo{
			Name:     "Redis Cloud",
			Homepage: sdk.URL("https://redis.com/"),
		},
		Credentials: []schema.CredentialType{
			RedisCloudAPIKey(),
		},
	}
}
