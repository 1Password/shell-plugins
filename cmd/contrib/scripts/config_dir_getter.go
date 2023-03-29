package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// Retrieve the config directory to read or write to
func main() {
	opConfigDir, _ := os.LookupEnv("OP_CONFIG_DIR")
	xdgConfigHome, _ := os.LookupEnv("XDG_CONFIG_HOME")
	home, _ := homedir.Dir()

	// This logic is based on the order of precedence outlined in the CLI documentation:
	// https://developer.1password.com/docs/cli/config-directories
	configDirPaths := []string{}
	if opConfigDir != "" {
		configDirPaths = append(configDirPaths, opConfigDir)
	}
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

	// If we reach this point then none of those directories exist (op is executed
	// for the first time). Default to the last entry in the list.
	if len(configDirPaths) > 0 {
		fmt.Print(configDirPaths[len(configDirPaths)-1])
	}
}
