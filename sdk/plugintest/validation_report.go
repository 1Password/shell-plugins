package plugintest

import (
	"fmt"
	"github.com/fatih/color"
	"os"

	"github.com/1Password/shell-plugins/sdk/schema"
)

func PrintValidationReport(plugin schema.Plugin) {
	reports := plugin.DeepValidate()
	printer := &ValidationReportPrinter{
		Reports: reports,
		Format:  PrintFormat{}.ValidationReportFormat(),
	}
	printer.Print()
}

func PrintValidateAllReport(plugins []schema.Plugin) {
	reports := make(map[string][]schema.ValidationReport)

	for _, p := range plugins {
		errorReports := FilterErrorReports(p.DeepValidate())
		reports[p.Name] = errorReports
	}

	printer := &ErrorReportPrinter{
		Reports: reports,
		Format:  PrintFormat{}.ValidationReportFormat(),
	}
	printer.Print()
}

func FilterErrorReports(reports []schema.ValidationReport) []schema.ValidationReport {
	var errorReports []schema.ValidationReport

	for _, report := range reports {
		errReport := schema.ValidationReport{Heading: report.Heading, Checks: []schema.ValidationCheck{}}
		for _, check := range report.Checks {
			if !check.Assertion && check.Severity == schema.ValidationSeverityError {
				errReport.Checks = append(errReport.Checks, check)
			}
		}
		errorReports = append(errorReports, errReport)
	}

	return errorReports
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

type ErrorReportPrinter struct {
	Reports map[string][]schema.ValidationReport
	Format  PrintFormat
}

func (printer *ErrorReportPrinter) Print() {
	var shouldExitWithError bool
	if printer.Reports == nil || len(printer.Reports) == 0 {
		return
	}

	for pluginName, reports := range printer.Reports {
		for _, report := range reports {
			if report.HasErrors() {
				printer.Format.Heading.Printf("Plugin %s has errors:\n", pluginName)
				printer.printHeading(report.Heading)
				printer.printChecks(report.Checks)
			}
		}
	}

	if shouldExitWithError {
		os.Exit(1)
	}
}

func (printer *ErrorReportPrinter) printHeading(heading string) {
	printer.Format.Heading.Printf("# %s\n\n", heading)
}

func (printer *ErrorReportPrinter) printChecks(checks []schema.ValidationCheck) {
	for _, c := range checks {
		printer.Format.Error.Printf("✘ %s\n", c.Description)
	}
	fmt.Println()
}
