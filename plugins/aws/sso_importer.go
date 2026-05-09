package aws

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const (
	profileSectionPrefix    = "profile "
	ssoSessionSectionPrefix = "sso-session "
)

// TrySSOConfigFile looks for AWS IAM Identity Center profiles in ~/.aws/config.
// It supports both the legacy form (sso_start_url on the profile) and the
// consolidated form (sso_session = NAME referencing an [sso-session NAME] section).
func TrySSOConfigFile() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		sourcePath := os.Getenv("AWS_CONFIG_FILE")
		if sourcePath == "" {
			sourcePath = "~/.aws/config"
		}

		configPath := sourcePath
		if strings.HasPrefix(configPath, "~") {
			configPath = in.FromHomeDir(strings.TrimPrefix(configPath, "~"))
		} else {
			configPath = in.FromRootDir(configPath)
		}

		contents, err := os.ReadFile(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				return
			}
			attempt := out.NewAttempt(importer.SourceFile(sourcePath))
			attempt.AddError(err)
			return
		}

		if len(contents) == 0 {
			return
		}

		attempt := out.NewAttempt(importer.SourceFile(sourcePath))

		configFile, err := importer.FileContents(contents).ToINI()
		if err != nil {
			attempt.AddError(err)
			return
		}

		ssoSessions := collectSSOSessions(configFile)

		for _, section := range configFile.Sections() {
			profileName, ok := profileNameFromSection(section.Name())
			if !ok {
				continue
			}

			fields, err := buildSSOFields(profileName, section, ssoSessions)
			if err != nil {
				attempt.AddError(err)
				continue
			}
			if fields == nil {
				continue
			}

			attempt.AddCandidate(sdk.ImportCandidate{
				Fields:   fields,
				NameHint: importer.SanitizeNameHint(profileName),
			})
		}
	}
}

// profileNameFromSection returns the bare profile name for a [profile NAME] or
// [default] section, and false for any other section type.
func profileNameFromSection(sectionName string) (string, bool) {
	if sectionName == defaultProfileName {
		return defaultProfileName, true
	}
	if strings.HasPrefix(sectionName, profileSectionPrefix) {
		return strings.TrimPrefix(sectionName, profileSectionPrefix), true
	}
	return "", false
}

// collectSSOSessions indexes [sso-session NAME] sections by their bare name.
func collectSSOSessions(configFile *ini.File) map[string]*ini.Section {
	sessions := make(map[string]*ini.Section)
	for _, section := range configFile.Sections() {
		if strings.HasPrefix(section.Name(), ssoSessionSectionPrefix) {
			name := strings.TrimPrefix(section.Name(), ssoSessionSectionPrefix)
			sessions[name] = section
		}
	}
	return sessions
}

// buildSSOFields returns the candidate fields for a profile section, or (nil, nil)
// if the profile is not SSO-bearing or is missing required keys. A non-nil error
// indicates a malformed reference (e.g. unknown sso_session) and should be reported.
func buildSSOFields(profileName string, section *ini.Section, ssoSessions map[string]*ini.Section) (map[sdk.FieldName]string, error) {
	hasSSOSession := keyHasValue(section, "sso_session")
	hasLegacyStartURL := keyHasValue(section, "sso_start_url")
	if !hasSSOSession && !hasLegacyStartURL {
		return nil, nil
	}

	fields := make(map[sdk.FieldName]string)

	if hasSSOSession {
		sessionName := section.Key("sso_session").Value()
		sessionSection, ok := ssoSessions[sessionName]
		if !ok {
			return nil, fmt.Errorf("profile %q references unknown sso-session %q", profileName, sessionName)
		}
		startURL := valueOrEmpty(sessionSection, "sso_start_url")
		region := valueOrEmpty(sessionSection, "sso_region")
		if startURL == "" || region == "" {
			return nil, fmt.Errorf("sso-session %q is missing sso_start_url or sso_region", sessionName)
		}
		fields[fieldname.SSOStartURL] = startURL
		fields[fieldname.SSORegion] = region
		fields[fieldname.SSOSession] = sessionName
	} else {
		fields[fieldname.SSOStartURL] = section.Key("sso_start_url").Value()
		fields[fieldname.SSORegion] = valueOrEmpty(section, "sso_region")
	}

	if v := valueOrEmpty(section, "sso_account_id"); v != "" {
		fields[fieldname.SSOAccountID] = v
	}
	if v := valueOrEmpty(section, "sso_role_name"); v != "" {
		fields[fieldname.SSORoleName] = v
	}
	if v := valueOrEmpty(section, "region"); v != "" {
		fields[fieldname.DefaultRegion] = v
	}

	if fields[fieldname.SSOStartURL] == "" || fields[fieldname.SSORegion] == "" ||
		fields[fieldname.SSOAccountID] == "" || fields[fieldname.SSORoleName] == "" {
		return nil, nil
	}

	return fields, nil
}

func keyHasValue(section *ini.Section, key string) bool {
	return section.HasKey(key) && section.Key(key).Value() != ""
}

func valueOrEmpty(section *ini.Section, key string) string {
	if !section.HasKey(key) {
		return ""
	}
	return section.Key(key).Value()
}
