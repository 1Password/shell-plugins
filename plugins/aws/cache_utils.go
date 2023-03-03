package aws

import (
	"context"
	"fmt"
	"github.com/1Password/shell-plugins/sdk"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	mfaCacheKey        = "sts-mfa"
	assumeRoleCacheKey = "sts-assume-role"
)

type StsCacheWriter struct {
	awsCacheKey string
	cache       sdk.CacheOperations
}

func (c StsCacheWriter) persist(credentials aws.Credentials) error {
	return c.cache.Put(c.awsCacheKey, credentials, credentials.Expires)
}

func newStsCacheWriter(key string, cache sdk.CacheOperations) *StsCacheWriter {
	return &StsCacheWriter{
		awsCacheKey: key,
		cache:       cache,
	}
}

type StsCacheProvider struct {
	awsCacheKey string
	cache       sdk.CacheState
}

func (c StsCacheProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	var cached aws.Credentials
	if ok := c.cache.Get(c.awsCacheKey, &cached); ok {
		return cached, nil
	}

	return aws.Credentials{}, fmt.Errorf("did not find cached credentials")
}

func getRoleCacheKey(awsConfig *confighelpers.Config) string {
	return assumeRoleCacheKey + awsConfig.RoleARN
}
