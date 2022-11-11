package plugintest

import (
	"fmt"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidationReportPrinter_printFields(t *testing.T) {
	printer := &ValidationReportPrinter{}
	fields := []schema.ValidationCheck{
		{Description: "error", Assertion: false, Severity: schema.ValidationSeverityError},
		{Description: "success", Assertion: true},
		{Description: "warning", Assertion: false, Severity: schema.ValidationSeverityWarning},
	}
	printer.sortChecks(&fields)

	assert.Equal(t, "success", fields[0].Description, fmt.Sprint("first filed should be success"))
	assert.Equal(t, "warning", fields[1].Description, fmt.Sprint("second filed should be warning"))
	assert.Equal(t, "error", fields[2].Description, fmt.Sprint("third filed should be error"))
}
