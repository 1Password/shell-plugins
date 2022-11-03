package plugins

import (
	"testing"
)

func TestValidatePlugins(t *testing.T) {
	for _, p := range registry {
		_, errs := p.Validate()
		if errs != nil {
			t.Logf("The '%s' plugin has validation errors: ", p.Name)
			for _, err := range errs {
				t.Log("\t", err)
			}
			t.Fail()
		}
	}
}
