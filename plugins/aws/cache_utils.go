package aws

import (
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	mfaCacheKeyID        = "sts-mfa"
	assumeRoleCacheKeyID = "sts-assume-role"
)

// stsCacheWriter writes aws temp credentials to cache using the awsCacheKey
type stsCacheWriter struct {
	awsCacheKey string
	cache       sdk.CacheOperations
}

func (c stsCacheWriter) Put(credentials aws.Credentials) error {
	return c.cache.Put(c.awsCacheKey, credentials, credentials.Expires)
}

func NewSTSCacheWriter(key string, cache sdk.CacheOperations) stsCacheWriter {
	return stsCacheWriter{
		awsCacheKey: key,
		cache:       cache,
	}
}

func getRoleCacheKey(roleArn string, accessKeyID string) string {
	return fmt.Sprintf("%s|%s|%s", assumeRoleCacheKeyID, accessKeyID, roleArn)
}

func getMfaCacheKey(accessKeyID string) string {
	return fmt.Sprintf("%s|%s", assumeRoleCacheKeyID, accessKeyID)
}
