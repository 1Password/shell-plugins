package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeFixture(t *testing.T, root, path, content string) {
	t.Helper()
	name := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(name), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSelection(t *testing.T) {
	plugins := []string{"alpha", "beta"}
	dependencies := map[string][]string{
		"alpha": {"plugins/alpha", "plugins/shared", "sdk"},
		"beta":  {"plugins/beta", "sdk"},
	}
	cases := []struct {
		name          string
		changed, want []string
	}{
		{"direct change", []string{"plugins/alpha/plugin.go"}, []string{"alpha"}},
		{"fixture change", []string{"plugins/alpha/test-fixtures/with space\nand newline.json"}, []string{"alpha"}},
		{"imported plugin change", []string{"plugins/shared/provisioner.go"}, []string{"alpha"}},
		{"prefix lookalike", []string{"plugins/alpha-other/plugin.go"}, []string{}},
		{"unrelated plugin", []string{"plugins/gamma/plugin.go"}, []string{}},
		{"shared dependency", []string{"sdk/provisioner.go"}, []string{"alpha", "beta"}},
		{"rename affects both", []string{"plugins/alpha/old.go", "plugins/beta/new.go"}, []string{"alpha", "beta"}},
		{"empty change", nil, []string{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := selectPlugins(plugins, dependencies, c.changed); !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %v; want %v", got, c.want)
			}
		})
	}
	for _, path := range []string{"go.mod", "go.sum", "sdk/importer/file.go", "cmd/ci/integration-matrix/main.go", ".github/workflows/integration.yaml"} {
		if !sharedChange([]string{path}) {
			t.Errorf("shared change was not recognized: %s", path)
		}
	}
	if sharedChange([]string{"README.md", "sdk-other/file.go"}) {
		t.Fatal("unrelated paths were marked shared")
	}
}

func TestDiscovery(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{"plugins/beta/one_integration_test.go", "plugins/alpha/one_integration_test.go", "plugins/alpha/two_integration_test.go", "plugins/other/unit_test.go"} {
		writeFixture(t, root, path, "package example\n")
	}
	got, err := discoverPlugins(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []string{"alpha", "beta"}) {
		t.Fatal(got)
	}
	writeFixture(t, root, "plugins/invalid;name/test_integration_test.go", "package example\n")
	if _, err := discoverPlugins(root); err == nil {
		t.Fatal("invalid directory accepted")
	}
	empty, err := discoverPlugins(t.TempDir())
	if err != nil || empty == nil || len(empty) != 0 {
		t.Fatalf("empty result must encode as [], got %v, %v", empty, err)
	}
}

func TestGitDiffAndGoDependencies(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "go.mod", "module example.test/integration\n\ngo 1.24\n")
	writeFixture(t, root, "plugins/shared/shared.go", "package shared\nconst Value = 1\n")
	writeFixture(t, root, "plugins/alpha/test_integration_test.go", "package alpha\nimport (\"testing\"; \"example.test/integration/plugins/shared\")\nfunc TestIntegration(t *testing.T) { _ = shared.Value }\n")
	runGit := func(args ...string) {
		t.Helper()
		if _, err := commandOutput(root, "git", args...); err != nil {
			t.Fatal(err)
		}
	}
	runGit("init", "--quiet")
	commit := func() {
		runGit("add", ".")
		runGit("-c", "user.name=Integration Test", "-c", "user.email=integration@example.invalid", "-c", "commit.gpgsign=false", "commit", "--quiet", "-m", "fixture")
	}
	oldPath := "plugins/shared/old space\nand newline.txt"
	newPath := "plugins/shared/new.txt"
	writeFixture(t, root, oldPath, "fixture")
	commit()
	runGit("checkout", "--quiet", "-b", "integration-change")
	if err := os.Rename(filepath.Join(root, oldPath), filepath.Join(root, newPath)); err != nil {
		t.Fatal(err)
	}
	commit()
	writeFixture(t, root, "README.md", "A later commit changes only documentation.")
	commit()
	runGit("checkout", "--quiet", "-b", "base-fixture", "HEAD~2")
	runGit("-c", "user.name=Integration Test", "-c", "user.email=integration@example.invalid", "-c", "commit.gpgsign=false", "merge", "--quiet", "--no-ff", "-m", "merge fixture", "integration-change")
	changed, err := changedFiles(root, "HEAD^1", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) != 3 {
		t.Fatalf("merge comparison must include both rename paths and the later documentation change: %q", changed)
	}
	seen := make(map[string]bool)
	for _, path := range changed {
		seen[path] = true
	}
	if !seen[oldPath] || !seen[newPath] {
		t.Fatalf("NUL-delimited paths were corrupted: %q", changed)
	}
	dirs, err := dependencyDirectories(root, "alpha")
	if err != nil {
		t.Fatal(err)
	}
	selected := selectPlugins([]string{"alpha"}, map[string][]string{"alpha": dirs}, changed)
	if !reflect.DeepEqual(selected, []string{"alpha"}) {
		t.Fatalf("real test dependency was missed: %v", dirs)
	}
	if _, err := changedFiles(root, "does-not-exist", "HEAD"); err == nil {
		t.Fatal("missing base must fail, not silently skip tests")
	}
	if _, err := dependencyDirectories(root, "does-not-exist"); err == nil {
		t.Fatal("broken dependency lookup must fail")
	}
}
