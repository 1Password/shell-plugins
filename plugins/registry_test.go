package plugins

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/schema"
)

func TestValidatePlugins(t *testing.T) {
	for _, p := range registry {
		_, report := p.Validate()
		for _, c := range report.Checks {
			if !c.Assertion && c.Severity == schema.ValidationSeverityError {
				t.Logf("The '%s' plugin has validation errors: %s", p.Name, c.Description)
				t.Fail()
			}
		}
	}
}
