package binance

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestBinanceCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, BinanceCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no without args": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
		"no for --help": {
			Args:              []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"no for --version": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"no for v1 public market data command": {
			Args:              []string{"book", "bnbusdt"},
			ExpectedNeedsAuth: false,
		},
		"yes for v1 order command": {
			Args:              []string{"buy", "-s", "BNBUSDT", "-t", "LIMIT", "-q", "0.05", "-p", "350", "-f", "GTC"},
			ExpectedNeedsAuth: true,
		},
		"no for v2 profile command": {
			Args:              []string{"profile", "list"},
			ExpectedNeedsAuth: false,
		},
		"no for v2 completion command": {
			Args:              []string{"completion", "zsh"},
			ExpectedNeedsAuth: false,
		},
		"yes for v2 signed command": {
			Args:              []string{"spot", "get-account"},
			ExpectedNeedsAuth: true,
		},
		"yes for a profile named after an exempt command": {
			Args:              []string{"spot", "get-account", "--profile", "profile"},
			ExpectedNeedsAuth: true,
		},
	})
}
