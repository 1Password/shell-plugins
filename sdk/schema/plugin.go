package schema

import (
	"fmt"
	"net/url"
	"strings"
)

// Plugin provides the schema for a single shell plugin. A plugin focuses on a single platform
// and can provide one or more credential types and one or more executables.
type Plugin struct {
	// The name of the plugin package, e.g. "aws" or "github". Should be the same name as the Go package.
	Name string

	// Details about the platform that the plugin covers.
	Platform PlatformInfo

	// One or more specifications for the credential types the plugin offers.
	Credentials []CredentialType

	// One or more specifications for the executables the plugin offers.
	Executables []Executable
}

// PlatformInfo provides information on the platform of the shell plugin.
type PlatformInfo struct {
	// The display name of the platform, e.g. "AWS" or "GitHub".
	Name string

	// The full URL of the homepage of the platform.
	Homepage *url.URL

	// (Optional) The full URL to the logo of the platform, in SVG or PNG format.
	Logo *url.URL
}

func (p Plugin) Validate() (bool, ValidationReport) {
	report := ValidationReport{
		Heading: fmt.Sprintf("Plugin: %s", p.Name),
		Checks:  []ValidationCheck{},
	}

	report.AddCheck(ValidationCheck{
		Description: "Has name set",
		Assertion:   p.Name != "",
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Name only using lowercase characters or digits",
		Assertion:   ContainsLowercaseLettersOrDigits(p.Name),
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has platform name set",
		Assertion:   p.Platform.Name != "",
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has platform homepage URL set",
		Assertion:   p.Platform.Homepage != nil,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has platform logo in SVG or PNG format",
		Assertion:   checkPlatformLogoFormat(p),
		Severity:    ValidationSeverityWarning,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has a credential type or executable defined",
		Assertion:   len(p.Credentials) > 0 && len(p.Executables) > 0,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has no more than one credential type defined",
		Assertion:   len(p.Credentials) == 1,
		Severity:    ValidationSeverityError,
	})

	return report.IsValid(), report
}

func (p Plugin) MakePluginValidationReports() []ValidationReport {
	var reports []ValidationReport

	_, pluginReport := p.Validate()
	reports = append(reports, pluginReport)

	for _, cred := range p.Credentials {
		_, credReport := cred.Validate()
		reports = append(reports, credReport)
	}

	for _, exe := range p.Executables {
		_, exeReport := exe.Validate()
		reports = append(reports, exeReport)
	}

	return reports
}

func checkPlatformLogoFormat(p Plugin) bool {
	if p.Platform.Logo == nil {
		return false
	}

	logoUrl := p.Platform.Logo.String()
	if strings.HasSuffix(logoUrl, ".png") || strings.HasSuffix(logoUrl, ".svg") {
		return true
	}

	return false
}
