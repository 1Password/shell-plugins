package plugintest

import (
	"fmt"
	"github.com/fatih/color"

	"github.com/1Password/shell-plugins/sdk/schema"
)

func PrintValidationReport(plugin schema.Plugin) {
	reports := plugin.MakePluginValidationReports()
	printer := &ValidationReportPrinter{
		Reports: reports,
		Format:  PrintFormat{}.ValidationReportFormat(),
	}
	printer.Print()
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
	vrp.printChecks(report.Checks)
}

// sortChecks in the order ["success", "warning", "error"]
func (vrp ValidationReportPrinter) sortChecks(checks *[]schema.ValidationCheck) {
	var successChecks []schema.ValidationCheck
	var warningChecks []schema.ValidationCheck
	var erroneousChecks []schema.ValidationCheck

	for _, c := range *checks {
		if c.Assertion {
			successChecks = append(successChecks, c)
			continue
		}

		if c.Severity == schema.ValidationSeverityWarning {
			warningChecks = append(warningChecks, c)
			continue
		}

		erroneousChecks = append(erroneousChecks, c)
	}

	*checks = append(successChecks, warningChecks...)
	*checks = append(*checks, erroneousChecks...)
}

func (vrp ValidationReportPrinter) printChecks(checks *[]schema.ValidationCheck) {
	vrp.sortChecks(checks)
	for _, c := range *checks {
		vrp.printCheck(c)
	}
	fmt.Println()
}

func (vrp ValidationReportPrinter) printHeading(heading string) {
	vrp.Format.Heading.Printf("# %s\n\n", heading)
}

func (vrp ValidationReportPrinter) printCheck(check schema.ValidationCheck) {
	if check.Assertion {
		vrp.Format.Success.Printf("✔ %s\n", check.Description)
		return
	}

	if check.Severity == schema.ValidationSeverityWarning {
		vrp.Format.Warning.Printf("⚠ %s\n", check.Description)
		return
	}

	vrp.Format.Error.Printf("𝘹%s\n", check.Description)
}
