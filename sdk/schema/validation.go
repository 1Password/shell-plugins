package schema

import (
	"log"
	"regexp"
	"strings"
)

type ValidationReport struct {
	Heading string
	Checks  *[]ValidationCheck
}

func (vr ValidationReport) AddCheck(check ValidationCheck) {
	*vr.Checks = append(*vr.Checks, check)
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
		log.Printf("error checking regexp %s", err)
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
		log.Printf("error checking regexp %s", err)
		return false
	}
	return matched
}

func IsValidReport(report ValidationReport) bool {
	isValid := true

	for _, check := range *report.Checks {
		if !check.Assertion {
			isValid = false
			break
		}
	}

	return isValid
}
