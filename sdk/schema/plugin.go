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

func (p Plugin) Validate() (isValid bool, errors []error) {
	if p.Name == "" {
		errors = append(errors, ErrMissingRequiredField("name"))
	}

	if p.Platform.Name == "" {
		errors = append(errors, ErrMissingRequiredField("platform.name"))
	}

	if p.Platform.Homepage == nil {
		errors = append(errors, ErrMissingRequiredField("platform.homepage"))
	}

	if len(p.Credentials) == 0 && len(p.Executables) == 0 {
		errors = append(errors, ErrMissingOneOfRequiredFields("credentials", "executables"))
	}

	if len(p.Credentials) > 1 {
		errors = append(errors, ErrNotYetSupported("provisioning multiple credentials to an executable"))
	}

	for _, cred := range p.Credentials {
		_, credErrors := cred.Validate()
		errors = append(errors, credErrors...)
	}

	for _, exe := range p.Executables {
		_, exeErrors := exe.Validate()
		errors = append(errors, exeErrors...)
	}

	return len(errors) == 0, errors
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
)
