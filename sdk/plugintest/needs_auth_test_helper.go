package plugintest

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/stretchr/testify/assert"
)

type NeedsAuthCase struct {
	Args              []string
	ExpectedNeedsAuth bool
}

func TestNeedsAuth(t *testing.T, rule sdk.NeedsAuthentication, cases map[string]NeedsAuthCase) {
	t.Helper()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Helper()
			in := sdk.NeedsAuthenticationInput{
				CommandArgs: c.Args,
			}
			assert.Equal(t, c.ExpectedNeedsAuth, rule(in), name)
		})
	}
}
