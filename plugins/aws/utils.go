package aws

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/99designs/aws-vault/v7/vault"
	"gopkg.in/ini.v1"
)

func getConfigSectionByProfile(configFile *ini.File, profileName string) *ini.Section {
	for _, section := range configFile.Sections() {
		if profileName == "default" && section.Name() == "default" {
			return section
		}

		// handle profile sections
		if section.Name() == fmt.Sprintf("profile %s", profileName) {
			return section
		}
	}

	return nil
}

func ExecuteSilently(f func() (*vault.ConfigFile, error)) func() (*vault.ConfigFile, error) {
	return func() (*vault.ConfigFile, error) {
		log.SetOutput(io.Discard)
		vaultConfig, err := f()
		defer log.SetOutput(os.Stderr)
		return vaultConfig, err
	}
}
