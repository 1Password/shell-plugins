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

	// The full URL to the logo of the platform, in SVG or PNG format.
	Logo *url.URL
}

func (p Plugin) Validate() (bool, ValidationReport) {
	report := ValidationReport{Heading: fmt.Sprintf("Plugin: %s", p.Name)}
	isValid, fields := validate(p)
	report.Fields = fields

	return isValid, report
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

func (p Plugin) ValidationSchema() ValidationSchema {
	return ValidationSchema{
		Fields: []ValidationSchemaField{
			{
				ReportText: "Has name set",
				Errors:     []error{},
				Validate: func() []error {
					var errors []error
					if p.Name == "" {
						errors = append(errors, ErrMissingRequiredField("name"))
					}
					return errors
				},
			},
			{
				ReportText: "Name only using lowercase characters or digits",
				Errors:     []error{},
				Validate: func() []error {
					var errors []error
					// TODO: implement
					errors = append(errors, fmt.Errorf("not implemented"))
					return errors
				},
			},
			{
				ReportText: "Has platform name set",
				Errors:     []error{},
				Validate: func() []error {
					var errors []error
					if p.Platform.Name == "" {
						errors = append(errors, ErrMissingRequiredField("platform.name"))
					}
					return errors
				},
			},
			{
				ReportText: "Has platform homepage URL set",
				Errors:     []error{},
				Validate: func() []error {
					var errors []error
					if p.Platform.Homepage == nil {
						errors = append(errors, ErrMissingRequiredField("platform.homepage"))
					}
					return errors
				},
			},
			{
				ReportText: "Has a credential type or executable defined",
				Errors:     []error{},
				Validate: func() []error {
					var errors []error
					if len(p.Credentials) == 0 && len(p.Executables) == 0 {
						errors = append(errors, ErrMissingOneOfRequiredFields("credentials", "executables"))
					}
					if len(p.Credentials) > 1 {
						errors = append(errors, ErrNotYetSupported("provisioning multiple credentials to an executable"))
					}
					return errors
				},
			},
		},
	}
}

var (
	ErrMissingRequiredField = func(fieldName string) error {
		return fmt.Errorf("missing required field: %s", fieldName)
	}

	ErrMissingOneOfRequiredFields = func(fields ...string) error {
		return fmt.Errorf("required to specify at least one of: %s", strings.Join(fields, ", "))
	}

	ErrNotYetSupported = func(action string) error {
		return fmt.Errorf("%s is not yet supporterd", action)
	}

	ErrTitleCase = func(fieldName string) error {
		return fmt.Errorf("%s should be in title case", fieldName)
	}
)
