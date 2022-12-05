package plugintest

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/stretchr/testify/assert"
)

func TestValidationReportPrinter_printFields(t *testing.T) {
	printer := &ValidationReportPrinter{}
	checks := []schema.ValidationCheck{
		{Description: "error", Assertion: false, Severity: schema.ValidationSeverityError},
		{Description: "success", Assertion: true},
		{Description: "warning", Assertion: false, Severity: schema.ValidationSeverityWarning},
	}
	checks = printer.sortChecks(checks)
	assert.Equal(t, "success", checks[0].Description, "first check should be success")
	assert.Equal(t, "warning", checks[1].Description, "second check should be warning")
	assert.Equal(t, "error", checks[2].Description, "third check should be error")
}
