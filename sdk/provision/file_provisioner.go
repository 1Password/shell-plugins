package provision

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/1Password/shell-plugins/sdk"
)

// FileProvisioner provisions one or more secrets as a temporary file.
type FileProvisioner struct {
	sdk.Provisioner

	fileContents        ItemToFileContents
	outfileName         string
	outpathFixed        string
	outpathEnvVar       string
	outdirEnvVar        string
	argPlacementMode    ArgPlacementMode
	outpathArgTemplates []string
}

type ItemToFileContents func(in sdk.ProvisionInput) ([]byte, error)

type ArgPlacementMode int

const (
	Unset ArgPlacementMode = iota
	Prepend
	Append
)

// FieldAsFile can be used to store the value of a single field as a file.
func FieldAsFile(fieldName sdk.FieldName) ItemToFileContents {
	return ItemToFileContents(func(in sdk.ProvisionInput) ([]byte, error) {
		if value, ok := in.ItemFields[fieldName]; ok {
			return []byte(value), nil
		} else {
			return nil, fmt.Errorf("no value present in the item for field '%s'", fieldName)
		}
	})
}

// TempFile returns a file provisioner and takes a function that maps a 1Password item to the contents of
// a single file.
func TempFile(fileContents ItemToFileContents, opts ...FileOption) FileProvisioner {
	p := FileProvisioner{
		fileContents: fileContents,
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// FileOption can be used to influence the behavior of the file provisioner.
type FileOption func(*FileProvisioner)

// AtFixedPath can be used to tell the file provisioner to store the credential at a specific location, instead of
// an autogenerated temp dir. This is useful for executables that can only load credentials from a specific path.
func AtFixedPath(path string) FileOption {
	return func(p *FileProvisioner) {
		p.outpathFixed = path
	}
}

// Filename can be used to tell the file provisioner to store the credential with a specific name, instead of
// an autogenerated name. The specified filename will be appended to the path of the autogenerated temp dir.
// Gets ignored if the provision.AtFixedPath option is also set.
func Filename(name string) FileOption {
	return func(p *FileProvisioner) {
		p.outfileName = name
	}
}

// SetPathAsEnvVar can be used to provision the temporary file path as an environment variable.
func SetPathAsEnvVar(envVarName string) FileOption {
	return func(p *FileProvisioner) {
		p.outpathEnvVar = envVarName
	}
}

// SetOutputDirAsEnvVar can be used to provision the directory of the output file as an environment variable.
func SetOutputDirAsEnvVar(envVarName string) FileOption {
	return func(p *FileProvisioner) {
		p.outdirEnvVar = envVarName
	}
}

// AppendArgs appends arguments to the command line for a FileProvisioner.
// This is particularly useful when you need to add arguments that reference the output file path.
// The output path is available as "{{ .Path }}" within the provided argument templates.
// For example:
// * `AppendArgs("--log", "{{ .Path }}")` results in `--log /path/to/tempfile`.
// * `AppendArgs("--log={{ .Path }}")` results in `--log=/path/to/tempfile`.
func AppendArgs(argTemplates ...string) FileOption {
	return func(p *FileProvisioner) {
		p.argPlacementMode = Append
		p.outpathArgTemplates = argTemplates
	}
}

// PrependArgs prepends arguments to the command line for a FileProvisioner.
// This is particularly useful when you need to add arguments that reference the output file path.
// The output path is available as "{{ .Path }}" within the provided argument templates.
// For example:
// * `PrependArgs("--input", "{{ .Path }}")` results in `--input /path/to/tempfile`.
// * `PrependArgs("--input={{ .Path }}")` results in `--input=/path/to/tempfile`.
//
// The arguments provided are added before any pre-existing arguments in the command line, but after the command itself.
func PrependArgs(argTemplates ...string) FileOption {
	return func(p *FileProvisioner) {
		p.argPlacementMode = Prepend
		p.outpathArgTemplates = argTemplates
	}
}

func (p FileProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	contents, err := p.fileContents(in)
	if err != nil {
		out.AddError(err)
		return
	}

	outpath := ""
	if p.outpathFixed != "" {
		// Default to the provision.AtFixedPath option
		outpath = p.outpathFixed
	} else if p.outfileName != "" {
		// Fall back to the provision.Filename option
		outpath = in.FromTempDir(p.outfileName)
	} else {
		// If both are undefined, resort to generating a random filename
		fileName, err := randomFilename()
		if err != nil {
			// This should only fail in rare circumstances
			out.AddError(fmt.Errorf("generating random file name: %s", err))
			return
		}
		outpath = in.FromTempDir(fileName)
	}

	out.AddSecretFile(outpath, contents)

	if p.outpathEnvVar != "" {
		// Populate the specified environment variable with the output path.
		out.AddEnvVar(p.outpathEnvVar, outpath)
	}

	if p.outdirEnvVar != "" {
		// Populate the specified environment variable with the output dir.
		dir := filepath.Dir(outpath)
		out.AddEnvVar(p.outpathEnvVar, dir)
	}

	// Add args to specify the output path.
	if p.argPlacementMode != Unset {
		tmplData := struct{ Path string }{
			Path: outpath,
		}

		// Resolve arg templates with the resulting output path injected.
		// Example: "--config-file={{ .Path }}" => "--config-file=/tmp/file"
		argsResolved := make([]string, len(p.outpathArgTemplates))
		for i, tmplStr := range p.outpathArgTemplates {
			tmpl, err := template.New("arg").Parse(tmplStr)
			if err != nil {
				out.AddError(err)
				return
			}

			var result bytes.Buffer
			err = tmpl.Execute(&result, tmplData)
			if err != nil {
				out.AddError(err)
				return
			}

			argsResolved[i] = result.String()
		}

		switch p.argPlacementMode {
		case Append:
			out.AppendArgs(argsResolved...)
		case Prepend:
			out.PrependArgs(argsResolved...)
		default:
			out.AddError(fmt.Errorf("invalid argument placement mode"))
		}
	}
}

func (p FileProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: deleting the files gets taken care of.
}

func (p FileProvisioner) Description() string {
	return "Provision secret file"
}

func randomFilename() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
