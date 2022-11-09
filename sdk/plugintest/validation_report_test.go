package plugintest

import (
	"errors"
	"fmt"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidationReportPrinter_printFields(t *testing.T) {
	printer := &ValidationReportPrinter{}
	fields := []schema.ValidationReportField{
		{ReportText: "error", Optional: false, Errors: []error{errors.New("1")}},
		{ReportText: "success", Errors: []error{}},
		{ReportText: "warning", Optional: true, Errors: []error{errors.New("1")}},
	}
	printer.sortFields(&fields)

	assert.Equal(t, "success", fields[0].ReportText, fmt.Sprint("first filed should be success"))
	assert.Equal(t, "warning", fields[1].ReportText, fmt.Sprint("second filed should be warning"))
	assert.Equal(t, "error", fields[2].ReportText, fmt.Sprint("third filed should be error"))
}
