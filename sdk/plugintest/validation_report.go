package plugintest

import (
	"fmt"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/fatih/color"
)

var printer = &ValidationReportPrinter{
	Format: PrintFormat{}.ValidationReportFormat(),
}

func PrintValidationReport(plugin schema.Plugin) {
	printer.Reports = plugin.DeepValidate()
	printer.Print()
}

func PrintReportIfErrors(plugin schema.Plugin) (hasErrors bool) {
	pluginReports := plugin.DeepValidate()
	for _, report := range pluginReports {
		if report.HasErrors() {
			printer.Reports = pluginReports
			printer.Print()
			return true
		}
	}

	return false
}

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

func (p *ValidationReportPrinter) Print() {
	if p.Reports == nil || len(p.Reports) == 0 {
		color.Cyan("No reports to print")
		return
	}

	for _, report := range p.Reports {
		p.PrintReport(report)
	}
}

func (p *ValidationReportPrinter) PrintReport(report schema.ValidationReport) {
	p.printHeading(report.Heading)
	p.printChecks(report.Checks)
}

// sortChecks in the order ["success", "warning", "error"]
func (p *ValidationReportPrinter) sortChecks(checks []schema.ValidationCheck) []schema.ValidationCheck {
	var successChecks []schema.ValidationCheck
	var warningChecks []schema.ValidationCheck
	var errorChecks []schema.ValidationCheck

	for _, c := range checks {
		if c.Assertion {
			successChecks = append(successChecks, c)
			continue
		}

		if c.Severity == schema.ValidationSeverityWarning {
			warningChecks = append(warningChecks, c)
			continue
		}

		errorChecks = append(errorChecks, c)
	}

	result := append(successChecks, warningChecks...)
	result = append(result, errorChecks...)

	return result
}

func (p *ValidationReportPrinter) printChecks(checks []schema.ValidationCheck) {
	for _, c := range p.sortChecks(checks) {
		p.printCheck(c)
	}
	fmt.Println()
}

func (p *ValidationReportPrinter) printHeading(heading string) {
	p.Format.Heading.Printf("# %s\n\n", heading)
}

func (p *ValidationReportPrinter) printCheck(check schema.ValidationCheck) {
	if check.Assertion {
		p.Format.Success.Printf("✔ %s\n", check.Description)
		return
	}

	if check.Severity == schema.ValidationSeverityWarning {
		p.Format.Warning.Printf("⚠ %s\n", check.Description)
		return
	}

	p.Format.Error.Printf("✘ %s\n", check.Description)
}
