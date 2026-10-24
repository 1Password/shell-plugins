package snowflakecli

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestSnowNeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, SnowflakeCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no args":                {Args: []string{}, ExpectedNeedsAuth: false},
		"help":                   {Args: []string{"sql", "--help"}, ExpectedNeedsAuth: false},
		"version":                {Args: []string{"--version"}, ExpectedNeedsAuth: false},
		"sql":                    {Args: []string{"sql", "-q", "select 1"}, ExpectedNeedsAuth: true},
		"connection test":        {Args: []string{"connection", "test"}, ExpectedNeedsAuth: true},
		"explicit connection":    {Args: []string{"sql", "-q", "select 1", "-c", "prod"}, ExpectedNeedsAuth: false},
		"explicit connection 2":  {Args: []string{"sql", "--connection", "prod"}, ExpectedNeedsAuth: false},
		"temporary connection":   {Args: []string{"sql", "-x", "--account", "acc"}, ExpectedNeedsAuth: false},
		"custom config file":     {Args: []string{"--config-file", "config.toml", "sql"}, ExpectedNeedsAuth: false},
		"connection list":        {Args: []string{"connection", "list"}, ExpectedNeedsAuth: false},
		"connection add":         {Args: []string{"connection", "add"}, ExpectedNeedsAuth: false},
		"connection set-default": {Args: []string{"connection", "set-default", "prod"}, ExpectedNeedsAuth: false},
		"helpers":                {Args: []string{"helpers", "check-version"}, ExpectedNeedsAuth: false},
		"init":                   {Args: []string{"init", "my_project"}, ExpectedNeedsAuth: false},
		"app bundle":             {Args: []string{"app", "bundle"}, ExpectedNeedsAuth: false},
		"app deploy":             {Args: []string{"app", "deploy"}, ExpectedNeedsAuth: true},
		"auth oidc read-token":   {Args: []string{"auth", "oidc", "read-token"}, ExpectedNeedsAuth: false},
		"custom-image validate":  {Args: []string{"custom-image", "validate"}, ExpectedNeedsAuth: false},
	})
}
