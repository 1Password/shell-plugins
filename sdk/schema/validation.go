package schema

import (
	"log"
	"regexp"
	"strings"
)

type ValidationReportSection string

type Validator interface {
	Validate() (bool, ValidationReport)
	ValidationSchema() ValidationSchema
}

type ValidationSchema struct {
	Fields []ValidationSchemaField
}

type ValidationSchemaField struct {
	ReportText string
	Optional   bool
	Errors     []error
	Validate   func() []error
}

type ValidationReport struct {
	Heading string
	Fields  []ValidationReportField
}

type ValidationReportField struct {
	ReportText string
	Optional   bool
	Errors     []error
}

func validate(v Validator) (bool, []ValidationReportField) {
	isValid := true
	var reportFields []ValidationReportField
	schema := v.ValidationSchema()

	for _, f := range schema.Fields {
		reportField := ValidationReportField{
			ReportText: f.ReportText,
			Optional:   f.Optional,
			Errors:     []error{},
		}
		errors := f.Validate()
		if len(errors) > 0 {
			reportField.Errors = errors
			isValid = false
		}
		reportFields = append(reportFields, reportField)
	}

	return isValid, reportFields
}

func IsErroneousField(field ValidationReportField) bool {
	return len(field.Errors) > 0
}

func IsOptionalField(field ValidationReportField) bool {
	return field.Optional
}

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
