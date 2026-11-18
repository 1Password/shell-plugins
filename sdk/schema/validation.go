package schema

import (
	"regexp"
	"strings"
)

type ValidationReport struct {
	Heading string
	Checks  []ValidationCheck
}

func (vr *ValidationReport) AddCheck(check ValidationCheck) {
	vr.Checks = append(vr.Checks, check)
}

func (vr *ValidationReport) IsValid() bool {
	isValid := true

	for _, check := range vr.Checks {
		if !check.Assertion {
			isValid = false
			break
		}
	}

	return isValid
}

func (vr *ValidationReport) HasErrors() bool {
	for _, check := range vr.Checks {
		if !check.Assertion && check.Severity == ValidationSeverityError {
			return true
		}
	}

	return false
}

type ValidationCheck struct {
	// Description explains what we want to validate
	Description string
	// Assertion
	Assertion bool
	// Severity is "warning" for Optional fields that are not passed and "error" for Required fields
	Severity ValidationSeverity
}

type ValidationSeverity string

const (
	ValidationSeverityWarning ValidationSeverity = "warning"
	ValidationSeverityError   ValidationSeverity = "error"
)

func IsTitleCaseWord(word string) bool {
	words := strings.Split(word, " ")
	if len(words) > 1 {
		return false
	}
	matched, err := regexp.Match("[A-Z][^\\s]*", []byte(word))
	if err != nil {
		return false
	}

	return matched
}

func IsTitleCaseString(str string) bool {
	if str == "" {
		return false
	}

	words := strings.Split(str, " ")
	if len(words) == 1 {
		return IsTitleCaseWord(words[0])
	}

	isTitleCaseStr := true
	for _, word := range words {
		if !IsTitleCaseWord(word) {
			isTitleCaseStr = false
			break
		}
	}

	return isTitleCaseStr
}

func ContainsLowercaseLettersOrDigits(str string) bool {
	matched, err := regexp.Match("^[a-z0-9]+$", []byte(str))
	if err != nil {
		return false
	}
	return matched
}

func CredentialReferencesInCredentialList(plugin Plugin) bool {
	for _, executable := range plugin.Executables {
		for _, execCredential := range executable.Uses {
			if execCredential.Name != "" {
				found := false
				for _, credential := range plugin.Credentials {
					if execCredential.Name == credential.Name {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}
	}
	return true
}

func NoDuplicateCredentials(plugin Plugin) bool {
	var ids []string
	for _, credential := range plugin.Credentials {
		ids = append(ids, credential.Name.ID().String())
	}

	return IsStringSliceASet(ids)
}

// AllCredentialsUsedWhenMultipleDefined reports whether every credential type of a plugin that defines
// more than one can be chosen for at least one of the plugin's executables. A credential type that no
// executable uses usually means it was added to the plugin definition but not to its executables.
//
// A credential type counts as used when an executable names it, or when an executable lets the user
// select from all credential types. Plugins with a single credential type, or without executables,
// always pass.
//
// This cannot see the executables of other plugins, so it also reports a credential type that is only
// offered for other plugins to use. That is why the corresponding validation check is a warning.
func AllCredentialsUsedWhenMultipleDefined(plugin Plugin) bool {
	if len(plugin.Credentials) <= 1 || len(plugin.Executables) == 0 {
		return true
	}

	used := map[string]bool{}
	for _, executable := range plugin.Executables {
		for _, usage := range executable.Uses {
			if usage.SelectFrom != nil && usage.SelectFrom.IncludeAllCredentials {
				return true
			}
			if usage.Name != "" && (usage.Plugin == "" || usage.Plugin == plugin.Name) {
				used[usage.Name.ID().String()] = true
			}
		}
	}

	for _, credential := range plugin.Credentials {
		if !used[credential.Name.ID().String()] {
			return false
		}
	}
	return true
}

func AreCredentialUsagesUniquelyIdentifiable(executable Executable) bool {
	var usageIds []string
	for _, credentialUsage := range executable.Uses {
		usageIds = append(usageIds, credentialUsage.ID())
	}

	return IsStringSliceASet(usageIds)
}

func IsStringSliceASet(slice []string) bool {
	for i, s := range slice {
		if i == len(slice)-1 {
			break
		}
		for _, ss := range slice[i+1:] {
			if ss == s {
				return false
			}
		}
	}

	return true
}
