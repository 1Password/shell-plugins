package plugins

import (
	"testing"
)

func TestValidatePlugins(t *testing.T) {
	for _, p := range registry {
		_, report := p.Validate()
		for _, f := range report.Fields {
			if len(f.Errors) > 0 {
				t.Logf("The '%s' plugin has validation errors: ", p.Name)
				for _, err := range f.Errors {
					t.Log("\t", err)
				}
				t.Fail()
			}
		}
	}
}
