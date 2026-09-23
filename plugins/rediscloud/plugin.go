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
			Homepage: sdk.URL("https://redis.com/redis-enterprise-cloud/overview/"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
	}
}
