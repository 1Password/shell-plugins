package aws

import (
	"fmt"

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
