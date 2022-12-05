package sdk

import (
	"context"
	"encoding/json"
	"path/filepath"
	"time"
)

// Provisioner provides hooks before and after the plugin's executable runs to provision
// and deprovision secrets or other means of authentication required for the executable to run.
type Provisioner interface {
	// Description describes what this provisioner does.
	Description() string

	// Provision gets called before running the plugin's executable to provision the necessary fields
	// from the 1Password item in a way that the executable understands.
	Provision(ctx context.Context, input ProvisionInput, output *ProvisionOutput)

	// Deprovision gets called after the plugin's executable exits, so that the plugin can clean up and
	// wipe any sensitive material created in the provision phase.
	Deprovision(ctx context.Context, input DeprovisionInput, output *DeprovisionOutput)
}

// ProvisionInput contains info that provisioners can use to provision credentials.
type ProvisionInput struct {
	// HomeDir is the path to current user's home directory.
	HomeDir string

	// TempDir is the path to a temporary directory that the provisioner can use to add files to.
	// This directory will automatically be deleted after the executable exits.
	TempDir string

	// DryRun can be used to opt out
	DryRun bool

	// Cache can contain data that got added in the provision step from previous runs for this credential.
	Cache CacheState

	// ItemFields contains the field names and their corresponding (sensitive) values.
	ItemFields map[string]string
}

// DeprovisionInput contains info that provisioners can use to deprovision credentials.
type DeprovisionInput struct {
	HomeDir string
	TempDir string
	DryRun  bool
}

// ProvisionOutput contains the sensitive values that the Provisioner outputs.
type ProvisionOutput struct {
	// Environment can be used to provision credentials as environment variable. The result of this will be added to the executable's environment.
	// The expected mapping is: environment variable name to (possibly sensitive) value.
	Environment map[string]string

	// CommandLine can be used provision credentials as command-line args. The result of this will be the actual (possibly sensitive) command
	// line that will be executed.
	CommandLine []string

	// Files can be used to provision credentials as files. The result of this will be automatically written to disk and deleted when the executable
	// exits. The expected mapping is: absolute file path to (possibly sensitive) file contents.
	Files map[string]OutputFile

	// Cache can be used to make data generated in this provision step available to the provision step of consecutive runs for this credential.
	// The data added to the cache will be encrypted and stored locally on disk, so it can be used to store sensitive data. To access the cached
	// data from previous runs, use Cache on ProvisionInput.
	Cache CacheOperations

	// Diagnostics can be used to report errors.
	Diagnostics Diagnostics
}

type DeprovisionOutput struct {
	Diagnostics Diagnostics
}

// OutputFile contains the sensitive file info and contents that the provisioner outputs.
type OutputFile struct {
	Contents []byte
}

type CacheState map[string]CacheEntry

type CacheOperations struct {
	Puts    map[string]CacheEntry
	Removes []string
}

// CacheEntry contains data persisted between consecutive provision runs
type CacheEntry struct {
	Data      []byte
	ExpiresAt time.Time
}

// AddEnvVar adds an environment variable to the provision output.
func (out *ProvisionOutput) AddEnvVar(name string, value string) {
	out.Environment[name] = value
}

// AddArgs can be used to add additional arguments to the command line of the provision output.
func (out *ProvisionOutput) AddArgs(args ...string) {
	out.CommandLine = append(out.CommandLine, args...)
}

// AddSecretFile can be used to add a file containing secrets to the provision output.
func (out *ProvisionOutput) AddSecretFile(path string, contents []byte) {
	out.AddFile(path, OutputFile{
		Contents: contents,
	})
}

// AddNonSecretFile can be used to add a file that does not contain secrets to the provision output.
func (out *ProvisionOutput) AddNonSecretFile(path string, contents []byte) {
	out.AddFile(path, OutputFile{
		Contents: contents,
	})
}

// AddFile can be used to add a file to the provision output.
func (out *ProvisionOutput) AddFile(path string, file OutputFile) {
	out.Files[path] = file
}

// AddError can be used to report an error to the provision output. If the provision output contains one
// or more errors, provisioning is considered failed.
func (out *ProvisionOutput) AddError(err error) {
	out.Diagnostics.Errors = append(out.Diagnostics.Errors, Error{err.Error()})
}

// FromHomeDir returns a path with the user's home directory prepended.
func (in *ProvisionInput) FromHomeDir(path ...string) string {
	return filepath.Join(append([]string{in.HomeDir}, path...)...)
}

// FromTempDir returns a path with the current execution's temp directory prepended.
func (in *ProvisionInput) FromTempDir(path ...string) string {
	return filepath.Join(append([]string{in.TempDir}, path...)...)
}

// Get returns the cached value at the specified key if it exists. The data can be returned either as a []byte
// or unmarshaled as JSON.
func (c CacheState) Get(key string, out any) (ok bool) {
	entry, ok := c[key]
	if !ok {
		return false
	}

	data := entry.Data
	switch out.(type) {
	case []byte:
		out = data
	default:
		err := json.Unmarshal(data, out)
		if err != nil {
			return false
		}
	}

	return true
}

// Put puts data into the cache at the specified key and with the specified TTL, which will be applied to the provision step of
// all consecutive runs, until the TTL is met or Remove is called. The data will be stored as a []byte or marshaled as JSON.
func (c *CacheOperations) Put(key string, data any, expiresAt time.Time) error {
	var marshaled []byte
	var err error

	switch data.(type) {
	case []byte:
		marshaled = data.([]byte)
	default:
		marshaled, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	c.Puts[key] = CacheEntry{
		ExpiresAt: expiresAt,
		Data:      marshaled,
	}

	return nil
}

// Remove removes data from the cache at the specified key, which will be applied to the provision step of all consecutive runs.
func (c *CacheOperations) Remove(key string) {
	c.Removes = append(c.Removes, key)
}
