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
	report := ValidationReport{Heading: fmt.Sprintf("Credential: %s", c.Name)}
	isValid, fields := validate(c)
	report.Fields = fields

	return isValid, report
}

func (c CredentialType) ValidationSchema() ValidationSchema {
	return ValidationSchema{
		Fields: []ValidationSchemaField{
			{
				ReportText: "Has name set",
				Validate: func() []error {
					var errors []error
					if c.Name == "" {
						errors = append(errors, ErrMissingRequiredField("name"))
					}
					return errors
				},
			},
			{
				ReportText: "Name is using title case",
				Validate: func() []error {
					var errors []error
					isTitleCase, err := IsTitleCase(c.Name)
					if !isTitleCase || err != nil {
						errors = append(errors, ErrTitleCase("name"))
					}
					return errors
				},
			},
			{
				ReportText: "Has docs URL set",
				Optional:   true,
				Validate: func() []error {
					var errors []error
					if c.DocsURL == nil {
						errors = append(errors, ErrMissingRequiredField("docsURL"))
					}
					return errors
				},
			},
			{
				ReportText: "Has management URL set",
				Optional:   true,
				Validate: func() []error {
					var errors []error
					if c.ManagementURL == nil {
						errors = append(errors, ErrMissingRequiredField("managementURL"))
					}
					return errors
				},
			},
			{
				ReportText: "Has at least 1 field",
				Validate: func() []error {
					var errors []error
					if len(c.Fields) == 0 {
						errors = append(errors, ErrMissingRequiredField("fields"))
					}
					return errors
				},
			},
			{
				ReportText: "All fields have name set",
				Validate: func() []error {
					var errors []error
					for _, f := range c.Fields {
						if f.Name == "" {
							errors = append(errors, ErrMissingRequiredField("field.name"))
						}
					}
					return errors
				},
			},
			{
				ReportText: "All field names are using title case",
				Validate: func() []error {
					var errors []error
					for _, f := range c.Fields {
						isTitleCase, err := IsTitleCase(f.Name)
						if !isTitleCase || err != nil {
							errors = append(errors, ErrTitleCase("field.name"))
						}
					}
					return errors
				},
			},
			{
				ReportText: "All fields have a description set",
				Validate: func() []error {
					var errors []error
					for _, f := range c.Fields {
						if f.MarkdownDescription == "" {
							errors = append(errors, ErrMissingRequiredField("field.markdownDescription"))
						}
					}
					return errors
				},
			},
			{
				ReportText: "All specified value compositions are valid",
				Validate: func() []error {
					var errors []error
					for _, f := range c.Fields {
						if f.MarkdownDescription == "" {
							errors = append(errors, ErrMissingRequiredField("field.name"))
						}
					}
					return errors
				},
			},
			{
				ReportText: "All specified value compositions are valid",
				Validate: func() []error {
					var errors []error
					for _, f := range c.Fields {
						comp := f.Composition
						if comp != nil {
							cs := comp.Charset
							if cs.Lowercase && cs.Uppercase && cs.Digits && cs.Symbols && len(cs.Specific) == 0 {
								errors = append(errors, ErrMissingOneOfRequiredFields(
									"composition.charset.lowercase",
									"composition.charset.uppercase",
									"composition.charset.digits",
									"composition.charset.symbols",
									"composition.charset.specific",
								))
							}
						}
					}
					return errors
				},
			},
			{
				ReportText: "Has at least 1 field that is secret",
				Validate: func() []error {
					var errors []error
					var credentialTypeHasSecret bool
					for _, f := range c.Fields {
						if f.Secret {
							credentialTypeHasSecret = true
							break
						}
					}
					if !credentialTypeHasSecret {
						errors = append(errors, fmt.Errorf("credential type must contain at least 1 secret field"))
					}
					return errors
				},
			},
			{
				ReportText: "Has a provisioner set",
				Validate: func() []error {
					var errors []error
					if c.Provisioner == nil {
						errors = append(errors, ErrMissingRequiredField("field.provisioner"))
					}
					return errors
				},
			},
			{
				ReportText: "Has an importer set",
				Optional:   true,
				Validate: func() []error {
					var errors []error
					if c.Importer == nil {
						errors = append(errors, ErrMissingRequiredField("field.importer"))
					}
					return errors
				},
			},
		},
	}
}
