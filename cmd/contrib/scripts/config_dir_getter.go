package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// Retrieve the config directory to read or write to
func main() {
	// This logic is based on the order of precedence outlined in the CLI documentation:
	// https://developer.1password.com/docs/cli/config-directories

	// If OP_CONFIG_DIR is set, use it immediately
	if opConfigDir, _ := os.LookupEnv("OP_CONFIG_DIR"); opConfigDir != "" {
		fmt.Print(opConfigDir)
		return
	}

	xdgConfigHome, _ := os.LookupEnv("XDG_CONFIG_HOME")
	home, _ := homedir.Dir()
	configDirPaths := []string{}
	if home != "" {
		// Legacy home
		configDirPaths = append(configDirPaths, filepath.Join(home, ".op"))
	}
	if xdgConfigHome != "" {
		// Legacy xdg
		configDirPaths = append(configDirPaths, filepath.Join(xdgConfigHome, ".op"))
	}
	if home != "" {
		// New home
		configDirPaths = append(configDirPaths, filepath.Join(home, ".config", "op"))
	}
	if xdgConfigHome != "" {
		// New xdg
		configDirPaths = append(configDirPaths, filepath.Join(xdgConfigHome, "op"))
	}

	for _, configDir := range configDirPaths {
		fileInfo, err := os.Stat(configDir)
		if err == nil && fileInfo.IsDir() {
			fmt.Print(configDir)
			return
		}
	}

	// None of those directories exist and OP_CONFIG_DIR is not set
	// Default to the last entry in the list
	if len(configDirPaths) > 0 {
		fmt.Print(configDirPaths[len(configDirPaths)-1])
	}
}
