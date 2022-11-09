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

	for _, f := range report.Fields {
		assert.Equal(t, true, isInvalidField(f), fmt.Sprintf("\"%s\" validation is erroneous", f.ReportText))
	}
}

func isInvalidField(field ValidationReportField) bool {
	return len(field.Errors) > 0
}
