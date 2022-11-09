package schema

type ValidationReportSection string

const (
	pluginSection      ValidationReportSection = "plugin"
	credentialsSection ValidationReportSection = "credentials"
	executablesSection ValidationReportSection = "executables"
)

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
