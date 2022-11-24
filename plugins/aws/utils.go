package aws

import (
	"gopkg.in/ini.v1"
	"strings"
)

func getConfigSectionByProfile(configFile *ini.File, profileName string) *ini.Section {
	for _, section := range configFile.Sections() {
		// handle [default] section
		if profileName == "default" && strings.Contains(section.Name(), "default") {
			return section
		}

		// handle [profile <profileName>] section
		if strings.Contains(section.Name(), "profile") && strings.Contains(section.Name(), profileName) {
			return section
		}
	}

	return nil
}
