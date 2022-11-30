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
	Fields    map[string]string
	NameHint  string
	ExpiresAt *time.Time
}

func (c *ImportCandidate) Equal(other ImportCandidate) bool {
	if len(c.Fields) != len(other.Fields) {
		return false
	}

	for key, value := range c.Fields {
		if value != other.Fields[key] {
			return false
		}
	}

	return true
}

type ImportSource struct {
	Env   []string
	Files []string
}

type ImportInput struct {
	HomeDir string
	RootDir string
}

type ImportOutput struct {
	Attempts []*ImportAttempt
}

type ImportAttempt struct {
	Candidates  []ImportCandidate
	Source      ImportSource
	Diagnostics Diagnostics
}

func (out *ImportOutput) Errors() (errors []Error) {
	for _, attempt := range out.Attempts {
		errors = append(errors, attempt.Diagnostics.Errors...)
	}
	return
}

func (out *ImportOutput) AllCandidates() (candidates []ImportCandidate) {
	for _, attempt := range out.Attempts {
		candidates = append(candidates, attempt.Candidates...)
	}
	return
}

func (out *ImportOutput) NewAttempt(src ImportSource) *ImportAttempt {
	attempt := &ImportAttempt{
		Source: src,
	}
	out.Attempts = append(out.Attempts, attempt)
	return attempt
}

func (out *ImportAttempt) AddCandidate(candidate ImportCandidate) {
	out.Candidates = append(out.Candidates, candidate)
}

func (out *ImportAttempt) AddError(err error) {
	out.Diagnostics.Errors = append(out.Diagnostics.Errors, Error{err.Error()})
}

func (in *ImportInput) FromHomeDir(path ...string) string {
	return filepath.Join(append([]string{in.HomeDir}, path...)...)
}

func (in *ImportInput) FromRootDir(path ...string) string {
	return filepath.Join(append([]string{in.RootDir}, path...)...)
}
