package aws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	confighelpers "github.com/99designs/aws-vault/v7/vault"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ssoTestServer stands in for the AWS SSO endpoints. It refreshes access tokens and hands out role
// credentials, and records what it was asked.
type ssoTestServer struct {
	*httptest.Server
	refreshRequests int
	bearerTokens    []string
	// rejectRefresh makes the server refuse to refresh, as it does once the SSO session has ended.
	rejectRefresh bool
}

func newSSOTestServer(t *testing.T) *ssoTestServer {
	t.Helper()
	s := &ssoTestServer{}
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		s.refreshRequests++
		w.Header().Set("Content-Type", "application/json")
		if s.rejectRefresh {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid_grant","error_description":"the refresh token is no longer valid"}`))
			return
		}
		_, _ = w.Write([]byte(`{"accessToken":"refreshed-access-token","refreshToken":"next-refresh-token","expiresIn":3600,"tokenType":"Bearer"}`))
	})
	mux.HandleFunc("/federation/credentials", func(w http.ResponseWriter, r *http.Request) {
		s.bearerTokens = append(s.bearerTokens, r.Header.Get("x-amz-sso_bearer_token"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"roleCredentials": map[string]any{
			"accessKeyId":     "ASIAEXAMPLE",
			"secretAccessKey": "example-secret",
			"sessionToken":    "example-session-token",
			"expiration":      time.Now().Add(time.Hour).UnixMilli(),
		}})
	})
	s.Server = httptest.NewServer(mux)
	t.Cleanup(s.Close)

	// Point the real AWS clients at the test server, so that the whole client stack is exercised.
	original := ssoClientConfig
	ssoClientConfig = func(region string) aws.Config {
		return aws.Config{Region: region, BaseEndpoint: aws.String(s.URL)}
	}
	t.Cleanup(func() { ssoClientConfig = original })
	return s
}

// writeSSOTokenCache writes a cached SSO login for the given cache key into a temporary home directory
// and returns the path of the file.
func writeSSOTokenCache(t *testing.T, cacheKey string, token map[string]string) string {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	path, err := ssocreds.StandardCachedTokenFilepath(cacheKey)
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	contents, err := json.Marshal(token)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, contents, 0o600))
	return path
}

func newTestSSOProvider(awsConfig *confighelpers.Config) (aws.CredentialsProvider, sdk.CacheOperations) {
	outCache := sdk.CacheOperations{Puts: map[string]sdk.CacheEntry{}}
	factory := SSOCacheProviderFactory{InCache: sdk.CacheState{}, OutCache: outCache}
	return factory.NewSSORoleCredentialsProvider(awsConfig), outCache
}

func TestSSOProviderRefreshesExpiredSessionToken(t *testing.T) {
	server := newSSOTestServer(t)
	tokenPath := writeSSOTokenCache(t, "my-sso", map[string]string{
		"accessToken":           "expired-access-token",
		"expiresAt":             time.Now().Add(-time.Minute).UTC().Format(time.RFC3339),
		"refreshToken":          "refresh-token",
		"clientId":              "client-id",
		"clientSecret":          "client-secret",
		"registrationExpiresAt": time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
	})

	provider, outCache := newTestSSOProvider(&confighelpers.Config{
		SSOSession:   "my-sso",
		SSOStartURL:  "https://example.awsapps.com/start",
		SSORegion:    "us-east-1",
		SSOAccountID: "123456789012",
		SSORoleName:  "ExampleRole",
	})
	creds, err := provider.Retrieve(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "ASIAEXAMPLE", creds.AccessKeyID)
	assert.Equal(t, 1, server.refreshRequests)
	assert.Equal(t, []string{"refreshed-access-token"}, server.bearerTokens)
	assert.Len(t, outCache.Puts, 1, "the role credentials should be cached")

	// The refreshed login is written back, so that the AWS CLI and later runs can use it.
	contents, err := os.ReadFile(tokenPath)
	require.NoError(t, err)
	var stored map[string]any
	require.NoError(t, json.Unmarshal(contents, &stored))
	assert.Equal(t, "refreshed-access-token", stored["accessToken"])
	assert.Equal(t, "next-refresh-token", stored["refreshToken"])
	info, err := os.Stat(tokenPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestSSOProviderDoesNotRefreshLegacyToken(t *testing.T) {
	server := newSSOTestServer(t)
	startURL := "https://example.awsapps.com/start"
	writeSSOTokenCache(t, startURL, map[string]string{
		"accessToken": "expired-access-token",
		"expiresAt":   time.Now().Add(-time.Minute).UTC().Format(time.RFC3339),
	})

	provider, _ := newTestSSOProvider(&confighelpers.Config{
		SSOStartURL:  startURL,
		SSORegion:    "us-east-1",
		SSOAccountID: "123456789012",
		SSORoleName:  "ExampleRole",
	})
	_, err := provider.Retrieve(context.Background())

	var invalidToken *ssocreds.InvalidTokenError
	require.True(t, errors.As(err, &invalidToken), "expected an invalid token error, got %v", err)
	assert.Equal(t, 0, server.refreshRequests)
	assert.Empty(t, server.bearerTokens)
}

func TestSSOProviderReportsUnusableSessionLoginAsInvalidToken(t *testing.T) {
	sessionConfig := &confighelpers.Config{
		SSOSession:   "my-sso",
		SSOStartURL:  "https://example.awsapps.com/start",
		SSORegion:    "us-east-1",
		SSOAccountID: "123456789012",
		SSORoleName:  "ExampleRole",
	}

	t.Run("no login in the token cache", func(t *testing.T) {
		server := newSSOTestServer(t)
		t.Setenv("HOME", t.TempDir())

		provider, _ := newTestSSOProvider(sessionConfig)
		_, err := provider.Retrieve(context.Background())

		var invalidToken *ssocreds.InvalidTokenError
		require.True(t, errors.As(err, &invalidToken), "expected an invalid token error, got %v", err)
		assert.Empty(t, server.bearerTokens)
	})

	t.Run("refresh is rejected", func(t *testing.T) {
		server := newSSOTestServer(t)
		server.rejectRefresh = true
		writeSSOTokenCache(t, "my-sso", map[string]string{
			"accessToken":           "expired-access-token",
			"expiresAt":             time.Now().Add(-time.Minute).UTC().Format(time.RFC3339),
			"refreshToken":          "refresh-token",
			"clientId":              "client-id",
			"clientSecret":          "client-secret",
			"registrationExpiresAt": time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
		})

		provider, _ := newTestSSOProvider(sessionConfig)
		_, err := provider.Retrieve(context.Background())

		var invalidToken *ssocreds.InvalidTokenError
		require.True(t, errors.As(err, &invalidToken), "expected an invalid token error, got %v", err)
		assert.Empty(t, server.bearerTokens)
	})
}

func TestSSOProviderUsesCachedRoleCredentialsOnlyWhileLoggedIn(t *testing.T) {
	startURL := "https://example.awsapps.com/start"
	legacyConfig := &confighelpers.Config{
		SSOStartURL:  startURL,
		SSORegion:    "us-east-1",
		SSOAccountID: "123456789012",
		SSORoleName:  "ExampleRole",
	}

	// newProviderWithCachedRoleCredentials returns a provider for which the plugin's cache already holds
	// role credentials that have not expired yet.
	newProviderWithCachedRoleCredentials := func(t *testing.T) aws.CredentialsProvider {
		t.Helper()
		cached := aws.Credentials{
			AccessKeyID:     "ASIACACHED",
			SecretAccessKey: "cached-secret",
			SessionToken:    "cached-session-token",
			CanExpire:       true,
			Expires:         time.Now().Add(30 * time.Minute),
		}
		data, err := json.Marshal(cached)
		require.NoError(t, err)
		cacheKey := getSSORoleCacheKey(legacyConfig.SSOAccountID, legacyConfig.SSORoleName, ssoSessionKey(legacyConfig))
		factory := SSOCacheProviderFactory{
			InCache:  sdk.CacheState{cacheKey: sdk.CacheEntry{Data: data, ExpiresAt: cached.Expires}},
			OutCache: sdk.CacheOperations{Puts: map[string]sdk.CacheEntry{}},
		}
		return factory.NewSSORoleCredentialsProvider(legacyConfig)
	}

	t.Run("login is still there", func(t *testing.T) {
		server := newSSOTestServer(t)
		// The access token has expired, which does not matter for role credentials that were already issued.
		writeSSOTokenCache(t, startURL, map[string]string{
			"accessToken": "expired-access-token",
			"expiresAt":   time.Now().Add(-time.Minute).UTC().Format(time.RFC3339),
		})

		creds, err := newProviderWithCachedRoleCredentials(t).Retrieve(context.Background())

		require.NoError(t, err)
		assert.Equal(t, "ASIACACHED", creds.AccessKeyID)
		assert.Empty(t, server.bearerTokens)
	})

	t.Run("login was removed by aws sso logout", func(t *testing.T) {
		server := newSSOTestServer(t)
		t.Setenv("HOME", t.TempDir())

		_, err := newProviderWithCachedRoleCredentials(t).Retrieve(context.Background())

		var invalidToken *ssocreds.InvalidTokenError
		require.True(t, errors.As(err, &invalidToken), "expected an invalid token error, got %v", err)
		assert.Empty(t, server.bearerTokens)
	})
}
