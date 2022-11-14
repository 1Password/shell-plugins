package schema

import (
	"fmt"
	"net/url"

	"github.com/1Password/shell-plugins/sdk"
)

// CredentialType provides the schema of a credential type that the plugin provides.
type CredentialType struct {
	// What the credential is called within the platform, e.g. "API Key" or "Personal Access Token".
	Name string

	// The field(s) on this credential type.
	Fields []CredentialField

	// (Optional) A URL to the documentation about this credential type.
	DocsURL *url.URL

	// (Optional) A URL to the dashboard, console, settings, etc. where this credential type can be created and revoked.
	ManagementURL *url.URL

	// (Optional) A function to scan the system for occurences of this credential type.
	Importer sdk.Importer

	// The provisioner to use provision and deprovision this credential to an exectuble.
	Provisioner sdk.Provisioner
}

// CredentialField provides the schema of a single field on a credential type.
type CredentialField struct {
	// The name of the field, e.g. "Token", "Password", or "Username".
	Name string

	// A description of the field.
	MarkdownDescription string

	// Whether this field is secret and should be concealed where possible.
	Secret bool

	// Whether this field is optional.
	Optional bool

	// (Optional) Describes how values of this field look like, such as the length, charset, etc.
	Composition *ValueComposition
}

func (c CredentialType) Field(name string) *CredentialField {
	for _, field := range c.Fields {
		if field.Name == name {
			return &field
		}
	}
	return nil
}

// ValueComposition describes what a value for a certain field looks like. This gets used for various purposes,
// including but not limited to the Save in 1Password functionality and secrets scanning functionality.
type ValueComposition struct {
	// The length of the value, if it's guaranteed to be of a fixed length.
	Length int

	// Which characters the value can consist of.
	Charset Charset

	// (Optional) A certain prefix that's always present on the value, as popularized by GitHub.
	Prefix string
}

type Charset struct {
	Uppercase bool
	Lowercase bool
	Digits    bool
	Symbols   bool
	Specific  []rune
}

func (c CredentialType) Validate() (bool, ValidationReport) {
	report := ValidationReport{
		Heading: fmt.Sprintf("Credential: %s", c.Name),
		Checks:  []ValidationCheck{},
	}

	isNameSet := c.Name != ""
	isDocUrlSet := c.DocsURL != nil
	isManagementUrlSet := c.ManagementURL != nil
	areAllFieldsHasNameSet := true
	areAllFieldsHasDescriptionSet := true
	areAllFieldsInTitleCase := true
	areAllCompositionsValid := true
	hasSecretField := false
	for _, f := range c.Fields {
		if f.Name == "" {
			areAllFieldsHasNameSet = false
		}
		if f.MarkdownDescription == "" {
			areAllFieldsHasDescriptionSet = false
		}
		if !IsTitleCaseString(f.Name) {
			areAllFieldsInTitleCase = false
		}
		comp := f.Composition
		if comp != nil {
			cs := comp.Charset
			if cs.Lowercase && cs.Uppercase && cs.Digits && cs.Symbols && len(cs.Specific) == 0 {
				areAllCompositionsValid = false
			}
		}
		if f.Secret {
			hasSecretField = true
		}
	}

	report.AddCheck(ValidationCheck{
		Description: "Has name set",
		Assertion:   isNameSet,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Name is using title case",
		Assertion:   IsTitleCaseString(c.Name),
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has documentation URL set",
		Assertion:   isDocUrlSet,
		Severity:    ValidationSeverityWarning,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has management URL set",
		Assertion:   isManagementUrlSet,
		Severity:    ValidationSeverityWarning,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has at least 1 field",
		Assertion:   len(c.Fields) > 0,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "All fields have name set",
		Assertion:   areAllFieldsHasNameSet,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "All field names are using title case",
		Assertion:   areAllFieldsInTitleCase,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "All fields have a description set",
		Assertion:   areAllFieldsHasDescriptionSet,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "All specified value compositions are valid",
		Assertion:   areAllCompositionsValid,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has at least 1 field that is secret",
		Assertion:   hasSecretField,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has a provisioner set",
		Assertion:   c.Provisioner != nil,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has an importer set",
		Assertion:   c.Importer != nil,
		Severity:    ValidationSeverityWarning,
	})

	return report.IsValid(), report
}
