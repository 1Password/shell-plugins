package plugintest

import (
	"fmt"
	"github.com/fatih/color"

	"github.com/1Password/shell-plugins/sdk/schema"
)

type PrintFormat struct {
	Heading *color.Color
	Warning *color.Color
	Error   *color.Color
	Success *color.Color
}

func (pf PrintFormat) ValidationReportFormat() PrintFormat {
	heading := color.New(color.FgCyan, color.Bold)
	warning := color.New(color.FgYellow)
	err := color.New(color.FgRed)
	success := color.New(color.FgGreen)

	return PrintFormat{
		Heading: heading,
		Warning: warning,
		Error:   err,
		Success: success,
	}
}

type ValidationReportPrinter struct {
	Reports []schema.ValidationReport
	Format  PrintFormat
}

func (vrp ValidationReportPrinter) Print() {
	if vrp.Reports == nil || len(vrp.Reports) == 0 {
		color.Cyan("No reports to print")
		return
	}

	for _, report := range vrp.Reports {
		vrp.PrintSectionReport(report)
	}
}

func (vrp ValidationReportPrinter) PrintSectionReport(report schema.ValidationReport) {
	vrp.printHeading(report.Heading)
	vrp.printFields(report.Fields)
}

func (vrp ValidationReportPrinter) sortFields(fields []schema.ValidationReportField) {
	for _, field := range fields {
		vrp.printField(field)
	}
	fmt.Println()
}

func (vrp ValidationReportPrinter) printFields(fields []schema.ValidationReportField) {
	for _, field := range fields {
		vrp.printField(field)
	}
	fmt.Println()
}

func (vrp ValidationReportPrinter) printHeading(heading string) {
	vrp.Format.Heading.Printf("# %s\n\n", heading)
}

func (vrp ValidationReportPrinter) printField(field schema.ValidationReportField) {
	if !schema.IsErroneousField(field) {
		vrp.Format.Success.Printf("✔ %s\n", field.ReportText)
		return
	}

	if schema.IsOptionalField(field) {
		vrp.Format.Warning.Printf("⚠ %s\n", field.ReportText)
		return
	}

	vrp.Format.Error.Printf("𝘹 %s\n", field.ReportText)
}

func PrintValidationReport(plugin schema.Plugin) {
	reports := plugin.MakePluginValidationReports()
	printer := &ValidationReportPrinter{
		Reports: reports,
		Format:  PrintFormat{}.ValidationReportFormat(),
	}
	printer.Print()
}
