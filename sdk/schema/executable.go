package schema

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
)

type Executable struct {
	// The entrypoint of the command that should be executed, e.g. ["aws"] or ["stripe"].
	Runs []string

	// The display name of the executable, e.g. "AWS CLI".
	Name string

	// Which credentials the executable requires to run and how these should be provisioned.
	Credentials []CredentialType

	// (Optional) A URL to the documentation about this executable.
	DocsURL *url.URL

	// (Optional) Whether the exectuable needs authentication for certain args.
	NeedsAuth sdk.NeedsAuthentication
}

func (e Executable) Validate() (bool, ValidationReport) {
	report := ValidationReport{Heading: fmt.Sprintf("Executable: %s", e.Name)}
	isValid, fields := validate(e)
	report.Fields = fields

	return isValid, report
}

func (e Executable) ValidationSchema() ValidationSchema {
	return ValidationSchema{
		Fields: []ValidationSchemaField{
			{
				ReportText: "Has name set",
				Validate: func() []error {
					var errors []error
					if e.Name == "" {
						errors = append(errors, ErrMissingRequiredField("name"))
					}
					return errors
				},
			},
			{
				ReportText: "Has docs URL set",
				Optional:   true,
				Validate: func() []error {
					var errors []error
					if e.DocsURL == nil {
						errors = append(errors, ErrMissingOptionalField("docsURL"))
					}
					return errors
				},
			},
			{
				ReportText: "Has specified which commands need authentication",
				Optional:   true,
				Validate: func() []error {
					var errors []error
					if e.NeedsAuth == nil {
						errors = append(errors, ErrMissingOptionalField("docsURL"))
					}
					return errors
				},
			},
			{
				ReportText: "Has executable command set",
				Validate: func() []error {
					var errors []error
					if len(e.Runs) == 0 {
						errors = append(errors, ErrMissingRequiredField("runs"))
					}
					return errors
				},
			},
			{
				ReportText: "Has a credential type defined",
				Validate: func() []error {
					var errors []error
					if len(e.Credentials) == 0 {
						errors = append(errors, ErrMissingRequiredField("credentials"))
					}
					return errors
				},
			},
		},
	}
}

func (e Executable) Command() string {
	return strings.Join(e.Runs, " ")
}

func CredentialWithProvisioner(credential CredentialType, provisioner sdk.Provisioner) CredentialType {
	credential.Provisioner = provisioner
	return credential
}
