package schema

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPluginValidateHasHeading(t *testing.T) {
	expectedHeading := "Plugin: test"
	p := Plugin{Name: "test"}
	_, report := p.Validate()

	assert.Equal(t, expectedHeading, report.Heading, fmt.Sprintf("plugin should have heading %s", expectedHeading))
}

func TestPluginValidateEachReportFieldHasError(t *testing.T) {
	p := Plugin{}
	_, report := p.Validate()

	for _, c := range *report.Checks {
		assert.Equal(t, false, c.Assertion, fmt.Sprintf("\"%s\" validation is erroneous", c.Description))
	}
}
