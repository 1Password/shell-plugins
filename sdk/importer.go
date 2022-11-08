package sdk

import (
	"context"
	"path/filepath"
	"time"
)

// Importer provides a hook for the plugin to scan the system for occurrences
// of a certain credential type, and returns every occurrence it can find.
type Importer func(ctx context.Context, in ImportInput, out *ImportOutput)

// ImportCandidate represents a single occurrence of a plugin's credential that was
// detected on the system.
type ImportCandidate struct {
	Fields    []ImportCandidateField
	NameHint  string
	ExpiresAt *time.Time
}

func (c ImportCandidate) Source() ImportSource {
	var source ImportSource
	for _, field := range c.Fields {
		source.Env = append(source.Env, field.Source.Env...)
		source.Files = append(source.Files, field.Source.Files...)
	}
	return source
}

func (c *ImportCandidate) Equal(other ImportCandidate) bool {
	if len(c.Fields) != len(other.Fields) {
		return false
	}
outer:
	for _, field := range c.Fields {
		for _, otherField := range other.Fields {
			if field.Equal(otherField) {
				continue outer
			}
		}
		return false
	}
	return true
}

// ImportCandidateField represents a single field and value of a credential type.
type ImportCandidateField struct {
	Field  string
	Value  string
	Source ImportSource
}

func (c *ImportCandidateField) Equal(other ImportCandidateField) bool {
	return c.Field == other.Field && c.Value == other.Value
}

type ImportSource struct {
	Env   []string
	Files []string
}

type ImportInput struct {
	HomeDir       string
	XDGConfigHome string
}

type ImportOutput struct {
	Candidates  []ImportCandidate
	Diagnostics Diagnostics
}

func (out *ImportOutput) AddCandidate(candidate ImportCandidate) {
	out.Candidates = append(out.Candidates, candidate)
}

func (out *ImportOutput) AddError(err error) {
	out.Diagnostics.Errors = append(out.Diagnostics.Errors, Error{err.Error()})
}

func (in *ImportInput) FromHomeDir(path ...string) string {
	return filepath.Join(append([]string{in.HomeDir}, path...)...)
}
