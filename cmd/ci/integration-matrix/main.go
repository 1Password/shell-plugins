// integration-matrix selects dependency-backed plugin tests for an optional CI job.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var pluginName = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

func main() {
	base := flag.String("base", "HEAD^1", "Base commit for the comparison")
	head := flag.String("head", "HEAD", "Head commit for the comparison")
	all := flag.Bool("all", false, "Select every plugin with integration tests")
	flag.Parse()
	root, err := os.Getwd()
	if err == nil {
		err = run(root, *base, *head, *all)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root, base, head string, all bool) error {
	plugins, err := discoverPlugins(root)
	if err != nil {
		return err
	}
	selected := plugins
	if !all && len(plugins) > 0 {
		changed, err := changedFiles(root, base, head)
		if err != nil {
			return err
		}
		if len(changed) == 0 {
			selected = []string{}
		} else if !sharedChange(changed) {
			dependencies := make(map[string][]string)
			for _, plugin := range plugins {
				directories, err := dependencyDirectories(root, plugin)
				if err != nil {
					return err
				}
				dependencies[plugin] = directories
			}
			selected = selectPlugins(plugins, dependencies, changed)
		}
	}
	return json.NewEncoder(os.Stdout).Encode(struct {
		Plugin []string `json:"plugin"`
	}{selected})
}

func discoverPlugins(root string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(root, "plugins", "*", "*_integration_test.go"))
	if err != nil {
		return nil, err
	}
	names := make(map[string]bool)
	for _, file := range files {
		name := filepath.Base(filepath.Dir(file))
		if !pluginName.MatchString(name) {
			return nil, fmt.Errorf("invalid integration plugin directory %q", name)
		}
		names[name] = true
	}
	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}
	sort.Strings(result)
	return result, nil
}

func commandOutput(root, name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = root
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w: %s", name, err, stderr.String())
	}
	return output, nil
}

func changedFiles(root, base, head string) ([]string, error) {
	revisions := make([]string, 0, 2)
	for _, ref := range []string{base, head} {
		output, err := commandOutput(root, "git", "rev-parse", "--verify", "--end-of-options", ref+"^{commit}")
		if err != nil {
			return nil, err
		}
		revisions = append(revisions, strings.TrimSpace(string(output)))
	}
	output, err := commandOutput(root, "git", "diff", "--name-only", "-z", "--no-renames", revisions[0], revisions[1], "--")
	if err != nil {
		return nil, err
	}
	if len(output) == 0 {
		return nil, nil
	}
	return strings.Split(strings.TrimSuffix(string(output), "\x00"), "\x00"), nil
}

func sharedChange(paths []string) bool {
	for _, path := range paths {
		if path == "go.mod" || path == "go.sum" || path == ".github/workflows/integration.yaml" ||
			strings.HasPrefix(path, "sdk/") || strings.HasPrefix(path, "cmd/ci/integration-matrix/") {
			return true
		}
	}
	return false
}

func dependencyDirectories(root, plugin string) ([]string, error) {
	output, err := commandOutput(root, "go", "list", "-deps", "-test", "-f", "{{.Dir}}", "./plugins/"+plugin)
	if err != nil {
		return nil, err
	}
	directories := make(map[string]bool)
	for _, dir := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if dir == "" {
			continue
		}
		relative, err := filepath.Rel(root, dir)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			continue
		}
		directories[filepath.ToSlash(relative)] = true
	}
	result := make([]string, 0, len(directories))
	for dir := range directories {
		result = append(result, dir)
	}
	sort.Strings(result)
	return result, nil
}

func selectPlugins(plugins []string, dependencies map[string][]string, changed []string) []string {
	selected := make([]string, 0)
	for _, plugin := range plugins {
		affected := false
		for _, directory := range dependencies[plugin] {
			for _, path := range changed {
				if directory == "." || strings.HasPrefix(path, directory+"/") {
					affected = true
				}
			}
		}
		if affected {
			selected = append(selected, plugin)
		}
	}
	return selected
}
