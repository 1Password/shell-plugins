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

func NoDuplicateCredentials(plugin Plugin) bool {
	var ids []string
	for _, credential := range plugin.Credentials {
		ids = append(ids, credential.Name.ID().String())
	}

	return IsStringSliceASet(ids)
}

func NoDuplicateCredentialUsages(executable Executable) bool {
	var usageIds []string
	for _, credentialUsage := range executable.Uses {
		usageIds = append(usageIds, credentialUsage.Name.ID().String()+"|"+credentialUsage.Plugin)
	}

	return IsStringSliceASet(usageIds)
}

func CredentialUsagesSpecifyCredentials(executable Executable) bool {
	for _, usage := range executable.Uses {
		if usage.Name == "" && usage.SelectFrom.ID == "" {
			return false
		}

		if usage.Name != "" && usage.SelectFrom.ID != "" {
			return false
		}
	}
	return true
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
