package backupdrill

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestBackupDrillCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, BackupDrillCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"yes for backup": {
			Args:              []string{"backup"},
			ExpectedNeedsAuth: true,
		},
		"yes for drill with a snapshot": {
			Args:              []string{"drill", "--snapshot", "2026-10-01T03-00-00-000Z"},
			ExpectedNeedsAuth: true,
		},
		"yes for restore dry run": {
			Args:              []string{"restore", "--database", "--dry-run"},
			ExpectedNeedsAuth: true,
		},
		"yes for a project named help": {
			Args:              []string{"backup", "--project-name", "help"},
			ExpectedNeedsAuth: true,
		},
		"no without args": {
			Args:              []string{},
			ExpectedNeedsAuth: false,
		},
		"no for --help": {
			Args:              []string{"backup", "--help"},
			ExpectedNeedsAuth: false,
		},
		"no for help command": {
			Args:              []string{"help", "drill"},
			ExpectedNeedsAuth: false,
		},
		"no for -h": {
			Args:              []string{"restore", "-h"},
			ExpectedNeedsAuth: false,
		},
		"no for -V": {
			Args:              []string{"-V"},
			ExpectedNeedsAuth: false,
		},
		"no for --version": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
	})
}
