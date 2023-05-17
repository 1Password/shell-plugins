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

	// (Optional) Whether the executable needs authentication for certain args.
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
	SelectFrom *CredentialSelection

	// (Optional) Whether the exectuable needs authentication for this credential. Works side by side with the executable's
	// `NeedsAuth`, which can still be used for more generic authentications opt-outs, such as the help flag.
	NeedsAuth sdk.NeedsAuthentication

	// Whether this credential is needed for the executable to run. If set to true, the executable cannot run without provisioning this credential.
	Optional bool
}

type CredentialSelection struct {
	// ID helps identify credentials chosen in this selection. This must be unique in relation to other selections specified in usages within its executable.
	ID string
	// IncludeAllCredentials specifies whether all credentials defined in all plugins should be included in the selection prompt when configuring this credential use.
	IncludeAllCredentials bool
	// AllowMultiple specifies whether multiple credentials can be selected to be part of this credential use.
	AllowMultiple bool
}

type CredentialReference struct {
	Name        sdk.CredentialName
	Plugin      string
	Provisioner sdk.Provisioner
}

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

	report.AddCheck(ValidationCheck{
		Description: "Credential Usages are uniquely identifiable inside an executable",
		Assertion:   CredentialUsagesUniquelyIdentifiable(e),
		Severity:    ValidationSeverityError,
	})

	return report.IsValid(), report
}

func (e Executable) Command() string {
	return strings.Join(e.Runs, " ")
}

func (c CredentialUsage) Validate() (bool, ValidationReport) {
	report := ValidationReport{
		Checks: []ValidationCheck{},
	}
	credRef, err := c.GetCredentialReference()
	report.AddCheck(ValidationCheck{
		Description: "If defined, a credential reference must have at least a Name",
		Assertion:   err == nil,
		Severity:    ValidationSeverityError,
	})

	selection, err := c.GetCredentialSelection()
	report.AddCheck(ValidationCheck{
		Description: "If defined, a credential selection must have an ID",
		Assertion:   err == nil,
		Severity:    ValidationSeverityError,
	})

	id, err := c.ID()
	report.AddCheck(ValidationCheck{
		Description: "Credential usage has either a credential reference or selection properly defined, but not both",
		Assertion:   err == nil && !(credRef != nil && selection != nil),
		Severity:    ValidationSeverityError,
	})
	report.Heading = fmt.Sprintf("Credential usage %s", id)

	return report.IsValid(), report
}

func (c CredentialUsage) GetCredentialReference() (*CredentialReference, error) {
	if c.Name != "" {
		return &CredentialReference{
			Name:        c.Name,
			Plugin:      c.Plugin,
			Provisioner: c.Provisioner,
		}, nil
	} else if c.Plugin != "" || c.Provisioner != nil {
		return nil, fmt.Errorf("plugin or provisioner specified without a credential name")
	}

	return nil, nil
}

func (c CredentialUsage) GetCredentialSelection() (*CredentialSelection, error) {
	if c.SelectFrom != nil {
		if c.SelectFrom.ID != "" {
			return c.SelectFrom, nil
		} else {
			return nil, fmt.Errorf("credential selection specified without an ID")
		}
	}

	return nil, nil
}

func (c CredentialUsage) ID() (string, error) {
	reference, err := c.GetCredentialReference()
	if err != nil {
		return "", err
	}

	selection, err := c.GetCredentialSelection()
	if err != nil {
		return "", err
	}

	if reference != nil {
		return strings.Join([]string{reference.Plugin, reference.Name.String()}, "|"), nil
	}

	if selection != nil {
		return selection.ID, nil
	}

	return "", fmt.Errorf("credential usage does not have a valid identifier")
}
