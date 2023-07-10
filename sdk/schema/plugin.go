package schema

import (
	"encoding/json"
	"fmt"
	"net/url"
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
}

func (p Plugin) Validate() (bool, ValidationReport) {
	report := ValidationReport{
		Heading: fmt.Sprintf("Plugin: %s", p.Name),
		Checks:  []ValidationCheck{},
	}

	report.AddCheck(ValidationCheck{
		Description: "Has plugin name set",
		Assertion:   p.Name != "",
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Plugin name only using lowercase characters or digits",
		Assertion:   ContainsLowercaseLettersOrDigits(p.Name),
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Plugin name not longer than 20 characters",
		Assertion:   len(p.Name) <= 20,
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
		Description: "Has a credential type or executable defined",
		Assertion:   len(p.Credentials) > 0 || len(p.Executables) > 0,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has no more than one credential type defined. Plugins with multiple credential types are not supported yet",
		Assertion:   len(p.Credentials) <= 1,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Credentials referenced in executables are included in the same plugin definition",
		Assertion:   CredentialReferencesInCredentialList(p),
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Credentials are uniquely identifiable inside a plugin",
		Assertion:   NoDuplicateCredentials(p),
		Severity:    ValidationSeverityError,
	})

	return report.IsValid(), report
}

func (p Plugin) DeepValidate() []ValidationReport {
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

		for _, usage := range exe.Uses {
			_, usageReport := usage.Validate()
			usageReport.Heading = fmt.Sprintf("Executable %s: %s", exe.Name, usageReport.Heading)
			reports = append(reports, usageReport)
		}
	}

	return reports
}

func (p Plugin) MarshalJSON() ([]byte, error) {
	if len(p.Credentials) == 0 {
		return nil, nil
	}
	plugin := struct {
		PlatformName  string `json:"platform_name"`
		ManagementURL string `json:"management_url"`
	}{
		PlatformName: p.Platform.Name,
		ManagementURL: func() string {
			if p.Credentials[0].ManagementURL != nil {
				return p.Credentials[0].ManagementURL.String()
			} else {
				return ""
			}
		}(),
	}
	return json.Marshal(plugin)
}
