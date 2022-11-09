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

	//if e.Name == "" {
	//	errors = append(errors, ErrMissingRequiredField("name"))
	//}
	//
	//if len(e.Runs) == 0 {
	//	errors = append(errors, ErrMissingRequiredField("runs"))
	//}
	//
	//return len(errors) == 0, errors
}

func (e Executable) ValidationSchema() ValidationSchema {
	return ValidationSchema{
		Fields: []ValidationSchemaField{
			{
				ReportText: "Has name set",
				Errors:     []error{},
				Validate: func() []error {
					var errors []error
					if e.Name == "" {
						errors = append(errors, ErrMissingRequiredField("name"))
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
