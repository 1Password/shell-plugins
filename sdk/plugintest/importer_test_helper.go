package plugintest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/stretchr/testify/assert"
)

// TestImporter will run the importer for each specified case, mounting the specified files in a temp dir,
// setting the environment variables, and configuring the home path using the temp dir.
func TestImporter(t *testing.T, importer sdk.Importer, cases map[string]ImportCase) {
	t.Helper()
	for description, c := range cases {
		if c.ExpectedOutput != nil && len(c.ExpectedCandidates) > 0 {
			t.Fatal("ExpectedOutput and ExpectedCandidates can't both be set in the same test case")
		}

		for envVarName, value := range c.Environment {
			t.Setenv(envVarName, value)
		}

		fsRoot := t.TempDir()
		in := sdk.ImportInput{
			HomeDir: filepath.Join(fsRoot, "~"),
		}

		for path, contents := range c.Files {
			path = filepath.Join(fsRoot, path)
			err := os.MkdirAll(filepath.Dir(path), 0700)
			if err != nil {
				t.Fatal(err)
			}

			err = os.WriteFile(path, []byte(contents), 0600)
			if err != nil {
				t.Fatal(err)
			}
		}

		ctx := context.Background()
		out := sdk.ImportOutput{}
		importer(ctx, in, &out)

		description = fmt.Sprintf("Import: %s", description)
		if c.ExpectedOutput != nil {
			assert.Equal(t, *c.ExpectedOutput, out, description)
		} else {
			assert.ElementsMatch(t, c.ExpectedCandidates, out.AllCandidates(), description)
		}

		for envVarName := range c.Environment {
			t.Setenv(envVarName, "")
		}
	}
}

type ImportCase struct {
	// Environment can be used to set environment variables for the importer test.
	Environment map[string]string

	// Files can be used to set files for the importer test, using the format: path -> contents.
	// For example: ~/.config/my-pugin/config -> '{"foo":"bar"}'. This is useful in conjunction with the
	// LoadFixture helper.
	Files map[string]string

	// ExpectedCandidates is a shorthand to set the expected import candidates. Mutually exclusive with ExpectedOutput.
	ExpectedCandidates []sdk.ImportCandidate

	// ExpectedOutput can be used to set the exact expected import output. Mutually exclusive with ExpectedCandidates.
	ExpectedOutput *sdk.ImportOutput
}
