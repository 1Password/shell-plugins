package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/1Password/shell-plugins/plugins"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/AlecAivazis/survey/v2"
)

const exampleSecretsCommandSuffix = "example-secrets"
const validateCommandSuffix = "validate"
const existsCommandSuffix = "exists"

func main() {
	command := os.Args[1]
	isPlugin, pluginName, pluginCommand := isPluginCommand(command)
	if isPlugin {
		plugin, err := plugins.Get(pluginName)
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasSuffix(pluginCommand, exampleSecretsCommandSuffix) {
			example := generateSecretsExample(plugin)
			fmt.Printf("%s", example)
			return
		}

		if strings.HasSuffix(pluginCommand, validateCommandSuffix) {
			plugintest.PrintValidationReport(plugin)
			return
		}

		if strings.HasSuffix(pluginCommand, existsCommandSuffix) {
			return
		}
	}

	if command == validateCommandSuffix {
		var shouldExitWithError bool
		for _, plugin := range plugins.List() {
			isReportPrinted := plugintest.PrintReportIfErrors(plugin)
			if isReportPrinted {
				shouldExitWithError = true
			}
		}

		if shouldExitWithError {
			os.Exit(1)
		}
		return
	}

	if command == "registry" {
		err := generatePluginRegistry()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if command == "new-plugin" {
		err := newPlugin()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

func isPluginCommand(command string) (isPluginCommand bool, pluginName string, pluginCommand string) {
	chunks := strings.Split(command, "/")
	if len(chunks) < 2 {
		return false, "", ""
	}
	return true, chunks[0], chunks[1]
}

func newPlugin() error {
	var questionnaire = []*survey.Question{
		{
			Name:     "Name",
			Prompt:   &survey.Input{Message: `Plugin name (e.g. "aws" or "github") [required]`},
			Validate: survey.Required,
		},
		{
			Name:     "PlatformName",
			Prompt:   &survey.Input{Message: `Platform name (e.g. "AWS" or "GitHub") [required]`},
			Validate: survey.Required,
		},
		{
			Name:     "Executable",
			Prompt:   &survey.Input{Message: `Executable name (e.g. "aws" or "gh") [required]`},
			Validate: survey.Required,
		},
		{
			Name: "CredentialName",
			Prompt: &survey.Input{
				Message: `Name of the credential type (e.g. "Access Key" or "Personal Access Token")`,
				Suggest: func(input string) []string {
					var suggestions []string
					for _, name := range credname.ListAll() {
						if strings.Contains(strings.ToLower(name.String()), strings.ToLower(input)) {
							suggestions = append(suggestions, name.String())
						}
					}
					return suggestions
				},
			},
			Validate:  validateCredentialName,
			Transform: transformCredentialName,
		},
		{
			Name:   "ExampleCredential",
			Prompt: &survey.Input{Message: `Paste in an example credential`},
		},
	}

	result := struct {
		Name              string
		PlatformName      string
		Executable        string
		CredentialName    string
		ExampleCredential string

		// Derived
		PlatformNameUpperCamelCase   string
		ValueComposition             schema.ValueComposition
		FieldName                    string
		FieldNameUpperCamelCase      string
		CredentialEnvVarName         string
		IsNewCredentialName          bool
		CredentialNameUpperCamelCase string
		CredentialNameSnakeCase      string
		TestCredentialExample        string
	}{}

	err := survey.Ask(questionnaire, &result)
	if err != nil {
		return err
	}

	if result.ExampleCredential != "" {
		result.ValueComposition = getValueComposition(result.ExampleCredential)
		result.TestCredentialExample = plugintest.ExampleSecretFromComposition(result.ValueComposition)
	} else {
		result.TestCredentialExample = plugintest.ExampleSecretFromComposition(schema.ValueComposition{
			Charset: schema.Charset{
				Uppercase: true,
				Lowercase: true,
				Digits:    true,
			},
			Length: 30,
		})
	}

	result.PlatformNameUpperCamelCase = strings.ReplaceAll(result.PlatformName, " ", "")

	credNameSplit := strings.Split(result.CredentialName, " ")

	result.CredentialNameUpperCamelCase = strings.Join(credNameSplit, "")
	result.CredentialNameSnakeCase = strings.ToLower(strings.Join(credNameSplit, "_"))

	result.IsNewCredentialName = true
	for _, existing := range credname.ListAll() {
		if result.CredentialName == existing.String() {
			result.IsNewCredentialName = false
			break
		}
	}

	// As a placeholder, assume the field name is the short version (max 7 chars) of the credential name, starting
	// from the last word. For example:
	// "Personal Access Token" => "Token"
	// "Secret Key" => "Key"
	// "API Key" => "API Key"
	var fieldNameSplit []string
	lengthCutoff := 7
	for i := range credNameSplit {
		word := credNameSplit[len(credNameSplit)-1-i]
		if len(strings.Join(append(fieldNameSplit, word), " ")) > lengthCutoff {
			break
		}

		fieldNameSplit = append([]string{word}, fieldNameSplit...)
	}
	result.FieldName = strings.Join(fieldNameSplit, " ")
	result.FieldNameUpperCamelCase = strings.Join(fieldNameSplit, "")
	result.CredentialEnvVarName = strings.ToUpper(strings.Join(append([]string{result.Name}, fieldNameSplit...), "_"))

	relativeDirPath := filepath.Join("plugins", result.Name)
	err = os.MkdirAll(relativeDirPath, 0777)
	if err != nil {
		return err
	}

	templates := []Template{pluginTemplate}
	if result.CredentialName != "" {
		templates = append(templates, credentialTemplate)
		templates = append(templates, credentialTestTemplate)
	}
	if result.Executable != "" {
		templates = append(templates, executableTemplate)
	}

	for _, tmpl := range templates {
		filenameTemplate, err := template.New("filename").Parse(tmpl.Filename)
		if err != nil {
			return err
		}

		var filenameBuf bytes.Buffer
		err = filenameTemplate.Execute(&filenameBuf, result)
		if err != nil {
			return err
		}
		filename := filenameBuf.String()

		contentsTemplate, err := template.New(filename).Parse(tmpl.Contents)
		if err != nil {
			return err
		}

		var contentsBuf bytes.Buffer
		err = contentsTemplate.Execute(&contentsBuf, result)
		if err != nil {
			return err
		}
		contents := contentsBuf.Bytes()

		err = os.WriteFile(filepath.Join(relativeDirPath, filename), contents, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func getValueComposition(value string) schema.ValueComposition {
	vc := schema.ValueComposition{
		Length: len(value),
	}

	for _, r := range value {
		if unicode.IsUpper(r) {
			vc.Charset.Uppercase = true
			continue
		}
		if unicode.IsLower(r) {
			vc.Charset.Lowercase = true
			continue
		}
		if unicode.IsDigit(r) {
			vc.Charset.Digits = true
			continue
		}
		if unicode.IsSymbol(r) {
			vc.Charset.Symbols = true
			continue
		}
	}

	vc.Prefix = getPossibleTokenPrefix(value)

	return vc
}

// getPossibleTokenPrefix tries to determine if the passed in value has a token prefix, as made popular by GitHub.
// Examples of such a prefix are `github_pat_`, `gph_`, `dop_v1_`, `glpat-`. See the code comments below for
// how that's determined.
func getPossibleTokenPrefix(value string) string {
	// If the value is shorter than 25 chars, it's unlikely to be a token that can contain a prefix.
	if len(value) < 20 {
		return ""
	}

	// A token prefix is likely to not be longer than 15 characters.
	substr := []rune(value)[:15]

	// A token prefix is likely to start with a lowercase char.
	if !unicode.IsLower(substr[0]) {
		return ""
	}

	// Trim all trailing chars until a delimiter is found, or else return an empty string.
	// A delimiter can be either an underscore (_) or a dash (-).
	return strings.TrimRightFunc(string(substr), func(r rune) bool {
		return r != '_' && r != '-'
	})
}

type Template struct {
	Filename string
	Contents string
}

var pluginTemplate = Template{
	Filename: "plugin.go",
	Contents: `package {{ .Name }}

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "{{ .Name }}",
		Platform: schema.PlatformInfo{
			Name:     "{{ .PlatformName }}",
			Homepage: sdk.URL("https://{{ .Name }}.com"), // TODO: Check if this is correct
		},
		{{- if .CredentialName }}
		Credentials: []schema.CredentialType{
			{{ .CredentialNameUpperCamelCase }}(),
		},
		{{- end }}
		{{- if .Executable }}
		Executables: []schema.Executable{
			{{ .PlatformNameUpperCamelCase }}CLI(),
		},
		{{- end }}
	}
}
`,
}

var credentialTemplate = Template{
	Filename: "{{ .CredentialNameSnakeCase }}.go",
	Contents: `package {{ .Name }}

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func {{ .CredentialNameUpperCamelCase }}() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.{{ .CredentialNameUpperCamelCase }},{{ if .IsNewCredentialName }} // TODO: Register name in project://sdk/schema/credname/names.go{{ end }}
		DocsURL:       sdk.URL("https://{{ .Name }}.com/docs/{{ .CredentialNameSnakeCase }}"), // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://console.{{ .Name }}.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.{{ .FieldNameUpperCamelCase }},
				MarkdownDescription: "{{ .FieldName }} used to authenticate to {{ .PlatformName }}.",
				Secret:              true,
				{{- if .ValueComposition.Length }}
				Composition: &schema.ValueComposition{
					{{- if .ValueComposition.Length }}
					Length: {{ .ValueComposition.Length }},
					{{- end }}
					{{- if .ValueComposition.Prefix }}
					Prefix: "{{ .ValueComposition.Prefix }}", // TODO: Check if this is correct
					{{- end }}
					Charset: schema.Charset{
						{{- if .ValueComposition.Charset.Uppercase }}
						Uppercase: true,
						{{- end }}
						{{- if .ValueComposition.Charset.Lowercase }}
						Lowercase: true,
						{{- end }}
						{{- if .ValueComposition.Charset.Digits }}
						Digits:    true,
						{{- end }}
						{{- if .ValueComposition.Charset.Symbols }}
						Symbols:   true,
						{{- end }}
					},
				},
				{{- end }}
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			Try{{ .PlatformNameUpperCamelCase }}ConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"{{ .CredentialEnvVarName }}": fieldname.{{ .FieldNameUpperCamelCase }}, // TODO: Check if this is correct
}

// TODO: Check if the platform stores the {{ .CredentialName }} in a local config file, and if so,
// implement the function below to add support for importing it.
func Try{{ .PlatformNameUpperCamelCase }}ConfigFile() sdk.Importer {
	return importer.TryFile("~/path/to/config/file.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config.{{ .FieldNameUpperCamelCase }} == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.{{ .FieldNameUpperCamelCase }}: config.{{ .FieldNameUpperCamelCase }},
		// 	},
		// })
	})
}

// TODO: Implement the config file schema
// type Config struct {
//	{{ .FieldNameUpperCamelCase }} string
// }
`,
}

var credentialTestTemplate = Template{
	Filename: "{{ .CredentialNameSnakeCase }}_test.go",
	Contents: `package {{ .Name }}

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func Test{{ .CredentialNameUpperCamelCase }}Provisioner(t *testing.T) {
	plugintest.TestProvisioner(t, {{ .CredentialNameUpperCamelCase }}().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.{{ .FieldNameUpperCamelCase }}: "{{ .TestCredentialExample }}",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"{{ .CredentialEnvVarName }}": "{{ .TestCredentialExample }}",
				},
			},
		},
	})
}

func Test{{ .CredentialNameUpperCamelCase }}Importer(t *testing.T) {
	plugintest.TestImporter(t, {{ .CredentialNameUpperCamelCase }}().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"{{ .CredentialEnvVarName }}": "{{ .TestCredentialExample }}",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.{{ .FieldNameUpperCamelCase }}: "{{ .TestCredentialExample }}",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in {{ .Name }}/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "{{ .TestCredentialExample }}",
			// 		},
			// 	},
			},
		},
	})
}
`,
}

var executableTemplate = Template{
	Filename: "{{ .Executable }}.go",
	Contents: `package {{ .Name }}

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func {{ .PlatformNameUpperCamelCase }}CLI() schema.Executable {
	return schema.Executable{
		Name:      "{{ .PlatformName }} CLI", // TODO: Check if this is correct
		Runs:      []string{"{{ .Executable }}"},
		DocsURL:   sdk.URL("https://{{ .Name }}.com/docs/cli"), // TODO: Replace with actual URL
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		{{- if .CredentialName }}
		Uses: []schema.CredentialUsage{
			{
				Name: credname.{{ .CredentialNameUpperCamelCase }},
			},
		},
		{{- end }}
	}
}
`,
}

func generateSecretsExample(plugin schema.Plugin) string {
	var example string

	for _, credential := range plugin.Credentials {
		example += credential.Name.String() + ":\n"
		for _, field := range credential.Fields {
			if field.Composition != nil {
				valueExample := plugintest.ExampleSecretFromComposition(*field.Composition)
				example += fmt.Sprintf("  %s: %s\n", field.Name, valueExample)
			}
		}
	}

	return example
}

func generatePluginRegistry() error {
	var plugins []string

	err := filepath.Walk("plugins", func(path string, info os.FileInfo, err error) error {
		if filepath.Base(path) == "plugin.go" {
			parent := filepath.Dir(path)

			if filepath.Base(filepath.Dir(parent)) == "plugins" {
				plugins = append(plugins, filepath.Base(parent))
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	tmpl, err := template.New("registry").Parse(pluginRegistryTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, plugins)
	if err != nil {
		return err
	}
	outfile := buf.Bytes()

	err = os.WriteFile(filepath.Join("plugins", "plugins.go"), outfile, 0600)
	if err != nil {
		return err
	}

	return nil
}

const pluginRegistryTemplate = `package plugins

// This file gets auto-generated by the "make registry" command, so should not be edited by hand.

import (
{{- range $plugin := . }}
	"github.com/1Password/shell-plugins/plugins/{{ $plugin }}"
{{- end }}
)

func init() {
{{- range $plugin := . }}
	Register({{ $plugin }}.New())
{{- end }}
}
`

func validateCredentialName(ans any) error {
	if str, ok := ans.(string); ok {
		if len(str) == 0 {
			return errors.New(`credential name must be titlecased, e.g. "Access Key" or "Personal Access Token"`)
		}

		words := strings.Split(str, " ")
		for _, word := range words {
			if unicode.IsLower(int32(word[0])) {
				return errors.New(`credential name must be titlecased, e.g. "Access Key" or "Personal Access Token"`)
			}
		}
	}

	return nil
}

func transformCredentialName(ans any) (newAns any) {
	if str, ok := ans.(string); ok {
		for _, name := range credname.ListAll() {
			if strings.EqualFold(name.String(), str) {
				return string(name)
			}
		}
	}

	return ans
}
