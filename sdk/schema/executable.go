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
	Uses []CredentialUsage

	// (Optional) A URL to the documentation about this executable.
	DocsURL *url.URL

	// (Optional) Whether the exectuable needs authentication for certain args.
	NeedsAuth sdk.NeedsAuthentication
}

type CredentialUsage struct {
	// (Optional) The name of the credential to use in the executable. Mutually exclusive with `SelectFrom`.
	Name sdk.CredentialName

	// (Optional) The plugin name that contains the credential. Defaults to the current package. This can be used to
	// include credentials from other plugins.
	Plugin string

	// (Optional) The provisioner to use to provision this credential to the executable. Overrides the DefaultProvisioner
	// set in the credential schema, so should only be used if this executable requires a custom configuration, that deviates
	// from the way the credential is usually provisioned.
	Provisioner sdk.Provisioner

	// (Optional) What this credential will be used for by the executable.
	Description string

	// (Optional) Instead of requiring a specific credential, have the user select from a list of compatible credentials.
	// Mutually exclusive with: `Name` and `Plugin`.
	SelectFrom CredentialSelection

	// (Optional) Whether the exectuable needs authentication for this credential. Works side by side with the executable's
	// `NeedsAuth`, which can still be used for more generic authentications opt-outs, such as the help flag.
	NeedsAuth sdk.NeedsAuthentication
}

type CredentialSelection struct {
	SelectionID    string
	PredefinedList []CredentialUsage
	Selector       CredentialSelector
	AllowMultiple  bool
}

type CredentialSelector string

func (e Executable) Validate() (bool, ValidationReport) {
	report := ValidationReport{
		Heading: fmt.Sprintf("Executable: %s", e.Name),
		Checks:  []ValidationCheck{},
	}

	report.AddCheck(ValidationCheck{
		Description: "Has name set",
		Assertion:   e.Name != "",
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has documentation URL set",
		Assertion:   e.DocsURL != nil,
		Severity:    ValidationSeverityWarning,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has specified which commands need authentication",
		Assertion:   e.NeedsAuth != nil,
		Severity:    ValidationSeverityWarning,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has executable command set",
		Assertion:   len(e.Runs) > 0,
		Severity:    ValidationSeverityError,
	})

	report.AddCheck(ValidationCheck{
		Description: "Has a credential type defined",
		Assertion:   len(e.Uses) > 0,
		Severity:    ValidationSeverityError,
	})

	return report.IsValid(), report
}

func (e Executable) Command() string {
	return strings.Join(e.Runs, " ")
}
