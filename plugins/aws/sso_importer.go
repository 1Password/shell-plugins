package aws

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"syscall"

	"gopkg.in/ini.v1"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const (
	profileSectionPrefix    = "profile "
	ssoSessionSectionPrefix = "sso-session "
)

// AWS account IDs are exactly 12 decimal digits; AWS regions are lowercase letters with one or two
// dashes separating segments and a trailing digit (e.g. `us-east-1`, `ap-southeast-2`).
var (
	ssoAccountIDRE = regexp.MustCompile(`^[0-9]{12}$`)
	ssoRegionRE    = regexp.MustCompile(`^[a-z]{2}-[a-z]+-[0-9]+$`)
)

// TrySSOConfigFile looks for AWS IAM Identity Center profiles in ~/.aws/config.
// It supports both the legacy form (sso_start_url on the profile) and the
// consolidated form (sso_session = NAME referencing an [sso-session NAME] section).
func TrySSOConfigFile() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		sourcePath := os.Getenv("AWS_CONFIG_FILE")
		envOverride := sourcePath != ""
		if !envOverride {
			sourcePath = "~/.aws/config"
		}

		// Match the SDK convention (sdk/importer/file_importer.go): only `~/` (with slash) is
		// expanded against the home directory. A bare `~root/...` form would otherwise be silently
		// joined to the *current* user's home, which is misleading at best and security-relevant if
		// an attacker can pre-place a file there.
		var configPath string
		switch {
		case strings.HasPrefix(sourcePath, "~/"):
			configPath = in.FromHomeDir(strings.TrimPrefix(sourcePath, "~/"))
		default:
			configPath = in.FromRootDir(sourcePath)
		}

		// When `AWS_CONFIG_FILE` is honoured, refuse anything that is not a regular file owned by
		// the current user. Defends against a malicious `.envrc` (direnv, asdf, mise) pinning the
		// importer to a hostile config file when `op` is run from that directory.
		if envOverride {
			if err := validateExternalConfigPath(configPath); err != nil {
				attempt := out.NewAttempt(importer.SourceFile(sourcePath))
				attempt.AddError(err)
				return
			}
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

		// Strict botocore-parity loader: `=` is the only key/value delimiter (botocore rejects `:`),
		// line continuation is disabled (botocore reads each `\` literally), and `Loose: true` lets
		// the loader produce a partial result instead of bricking the entire import on a single
		// malformed section. AllowShadows defaults to false, which matches botocore last-wins
		// semantics for duplicate `[profile X]` sections; we surface no diagnostic to keep parity.
		configFile, err := ini.LoadSources(ini.LoadOptions{
			KeyValueDelimiters: "=",
			IgnoreContinuation: true,
			Loose:              true,
		}, contents)
		if err != nil {
			attempt.AddError(err)
			// Fall through; with Loose: true the loader returns a partial *ini.File on most parse
			// errors and only short-circuits on truly catastrophic ones. Either way, valid sections
			// in the result still produce candidates.
		}
		if configFile == nil {
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

// validateExternalConfigPath enforces the safety properties for an `AWS_CONFIG_FILE` override:
// the path must resolve to a regular file owned by the current user. Anything else is rejected.
func validateExternalConfigPath(path string) error {
	fi, err := os.Lstat(path)
	if err != nil {
		// If the path does not exist, downstream `os.ReadFile` will report it; treat it as a
		// non-error here so we do not double-report.
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("AWS_CONFIG_FILE %q: %w", path, err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("AWS_CONFIG_FILE %q is a symlink; refusing to follow", path)
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("AWS_CONFIG_FILE %q is not a regular file", path)
	}
	if st, ok := fi.Sys().(*syscall.Stat_t); ok {
		if st.Uid != uint32(os.Geteuid()) {
			return fmt.Errorf("AWS_CONFIG_FILE %q is not owned by the current user", path)
		}
	}
	return nil
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

// buildSSOFields returns the candidate fields for a profile section, or (nil, nil) if the profile
// is not SSO-bearing or is missing required keys. A non-nil error indicates a malformed reference
// or an invalid value (bad URL, malformed account ID, etc.) and should be reported via
// `attempt.AddError`. Errors are scoped to a single profile so other valid profiles still import.
func buildSSOFields(profileName string, section *ini.Section, ssoSessions map[string]*ini.Section) (map[sdk.FieldName]string, error) {
	hasSSOSession := keyHasValue(section, "sso_session")
	hasLegacyStartURL := keyHasValue(section, "sso_start_url")
	if !hasSSOSession && !hasLegacyStartURL {
		return nil, nil
	}

	fields := make(map[sdk.FieldName]string)

	if hasSSOSession {
		sessionName := section.Key("sso_session").Value()
		if containsNUL(sessionName) {
			return nil, fmt.Errorf("profile %q sso_session value contains a NUL byte", profileName)
		}
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

	for _, v := range fields {
		if containsNUL(v) {
			return nil, fmt.Errorf("profile %q contains a NUL byte in one of its SSO fields", profileName)
		}
	}

	if err := validateSSOStartURL(fields[fieldname.SSOStartURL]); err != nil {
		return nil, fmt.Errorf("profile %q: %w", profileName, err)
	}
	if !ssoAccountIDRE.MatchString(fields[fieldname.SSOAccountID]) {
		return nil, fmt.Errorf("profile %q: sso_account_id %q is not a 12-digit AWS account ID", profileName, fields[fieldname.SSOAccountID])
	}
	if !ssoRegionRE.MatchString(fields[fieldname.SSORegion]) {
		return nil, fmt.Errorf("profile %q: sso_region %q is not a valid AWS region", profileName, fields[fieldname.SSORegion])
	}

	return fields, nil
}

// validateSSOStartURL enforces an HTTPS scheme and a non-empty host. We do not pin the host to
// `*.awsapps.com`: AWS supports custom SSO start URLs (e.g. `https://signin.aws.amazon.com/...`),
// and an over-tight allowlist would reject legitimate enterprise configurations.
func validateSSOStartURL(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("sso_start_url %q is not a valid URL: %w", s, err)
	}
	if u.Scheme != "https" {
		return fmt.Errorf("sso_start_url %q must use https://", s)
	}
	if u.Host == "" {
		return fmt.Errorf("sso_start_url %q has no host component", s)
	}
	return nil
}

func containsNUL(s string) bool {
	return strings.IndexByte(s, 0) >= 0
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
