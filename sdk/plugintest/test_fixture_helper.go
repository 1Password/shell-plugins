package plugintest

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// LoadFixture loads the test fixture file from the "test-fixtures" dir in the plugin directory.
// It fails the test if the file can't be loaded.
func LoadFixture(t *testing.T, filename string) string {
	t.Helper()

	_, testFilename, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatal()
	}

	fixturePath := filepath.Join(filepath.Dir(testFilename), "test-fixtures", filename)
	contents, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}

	return string(contents)
}
