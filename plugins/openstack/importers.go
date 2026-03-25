package openstack

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

// TryOpenStackCloudRC imports credentials from an OpenStack RC file (cloudrc / openrc.sh).
// These files are typically downloaded from the OpenStack Horizon dashboard and contain
// "export OS_*=value" shell statements that configure authentication environment variables.
func TryOpenStackCloudRC(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		fields := make(map[sdk.FieldName]string)

		scanner := bufio.NewScanner(strings.NewReader(contents.ToString()))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, "export ") {
				continue
			}
			line = strings.TrimPrefix(line, "export ")
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Strip surrounding quotes if present
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'') {
					value = value[1 : len(value)-1]
				}
			}
			if fieldName, ok := defaultEnvVarMapping[key]; ok && value != "" {
				fields[fieldName] = value
			}
		}

		if len(fields) > 0 {
			out.AddCandidate(sdk.ImportCandidate{Fields: fields})
		}
	})
}

// TryOpenStackCloudsYAMLFromEnvVar imports credentials from the clouds.yaml file pointed to
// by the OS_CLIENT_CONFIG_FILE environment variable, which the OpenStack CLI checks first
// before falling back to the standard locations. A secure.yaml in the same directory is
// also merged in if present.
func TryOpenStackCloudsYAMLFromEnvVar() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		cloudsPath := os.Getenv("OS_CLIENT_CONFIG_FILE")
		if cloudsPath == "" {
			return
		}
		// Derive secure.yaml as a sibling of the specified clouds file (unexpanded,
		// so TryOpenStackCloudsAndSecureYAML can expand it correctly).
		securePath := filepath.Join(filepath.Dir(cloudsPath), "secure.yaml")
		TryOpenStackCloudsAndSecureYAML(cloudsPath, securePath)(ctx, in, out)
	}
}

// TryOpenStackCloudsAndSecureYAML imports credentials by merging an OpenStack clouds.yaml
// file with the companion secure.yaml file at the same path. secure.yaml provides the base
// values (typically sensitive fields such as passwords) and clouds.yaml takes priority,
// following the OpenStack CLI convention for separating sensitive data from configuration.
func TryOpenStackCloudsAndSecureYAML(cloudsPath, securePath string) sdk.Importer {
	return importer.TryFile(cloudsPath, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var cloudsConfig cloudsYAML
		if err := contents.ToYAML(&cloudsConfig); err != nil {
			out.AddError(err)
			return
		}

		// Read secure.yaml (optional — silently skip if absent or no path given).
		var secureConfig cloudsYAML
		if securePath != "" {
			secureFilePath := expandPath(securePath, in)
			if secureBytes, err := os.ReadFile(secureFilePath); err == nil {
				_ = importer.FileContents(secureBytes).ToYAML(&secureConfig)
			}
		}

		// Build merged cloud map. clouds.yaml is the authoritative source of which clouds
		// exist — orphaned entries in secure.yaml (no matching clouds.yaml entry) are ignored.
		allClouds := make(map[string]cloudEntry)
		for name, entry := range cloudsConfig.Clouds {
			if secureEntry, exists := secureConfig.Clouds[name]; exists {
				allClouds[name] = mergeCloudEntry(secureEntry, entry)
			} else {
				allClouds[name] = entry
			}
		}

		cloudNames := make([]string, 0, len(allClouds))
		for name := range allClouds {
			cloudNames = append(cloudNames, name)
		}
		sort.Strings(cloudNames)

		for _, cloudName := range cloudNames {
			cloud := allClouds[cloudName]
			fields := make(map[sdk.FieldName]string)

			if cloud.Auth.AuthURL != "" {
				fields[fieldname.URL] = cloud.Auth.AuthURL
			}
			if cloud.Auth.Username != "" {
				fields[fieldname.Username] = cloud.Auth.Username
			}
			if cloud.Auth.Password != "" {
				fields[fieldname.Password] = cloud.Auth.Password
			}
			if cloud.Auth.ApplicationCredentialID != "" {
				fields[fieldApplicationCredentialID] = cloud.Auth.ApplicationCredentialID
			}
			if cloud.Auth.ApplicationCredentialSecret != "" {
				fields[fieldApplicationCredentialSecret] = cloud.Auth.ApplicationCredentialSecret
			}
			if cloud.Auth.ProjectName != "" {
				fields[fieldname.Project] = cloud.Auth.ProjectName
			}
			if cloud.Auth.ProjectID != "" {
				fields[fieldname.ProjectID] = cloud.Auth.ProjectID
			}
			if cloud.RegionName != "" {
				fields[fieldname.Region] = cloud.RegionName
			}
			if cloud.Auth.UserDomainName != "" {
				fields[fieldUserDomainName] = cloud.Auth.UserDomainName
			}
			// project_domain_name can appear inside auth or at the entry level
			if cloud.Auth.ProjectDomainName != "" {
				fields[fieldProjectDomainName] = cloud.Auth.ProjectDomainName
			} else if cloud.ProjectDomainName != "" {
				fields[fieldProjectDomainName] = cloud.ProjectDomainName
			}
			if cloud.Auth.ProjectDomainID != "" {
				fields[fieldProjectDomainID] = cloud.Auth.ProjectDomainID
			}
			if cloud.Interface != "" {
				fields[fieldInterface] = cloud.Interface
			}
			if cloud.IdentityAPIVersion != "" {
				fields[fieldIdentityAPIVersion] = cloud.IdentityAPIVersion
			}
			if cloud.AuthType != "" {
				fields[fieldAuthType] = cloud.AuthType
			}

			if len(fields) > 0 {
				out.AddCandidate(sdk.ImportCandidate{
					Fields:   fields,
					NameHint: importer.SanitizeNameHint(cloudName),
				})
			}
		}
	})
}

// TryOpenStackCloudsYAML imports credentials from an OpenStack clouds.yaml file only,
// without loading a companion secure.yaml. Use TryOpenStackCloudsAndSecureYAML when
// both files may be present.
func TryOpenStackCloudsYAML(path string) sdk.Importer {
	return TryOpenStackCloudsAndSecureYAML(path, "")
}

// expandPath resolves a path that may start with "~" to an absolute path using
// the home directory from the import context.
func expandPath(path string, in sdk.ImportInput) string {
	if strings.HasPrefix(path, "~") {
		return in.FromHomeDir(path[1:])
	}
	return in.FromRootDir(path)
}

// mergeCloudEntry merges two cloudEntry structs. The base provides default values;
// any non-empty field in override takes priority.
func mergeCloudEntry(base, override cloudEntry) cloudEntry {
	result := base
	result.Auth = mergeCloudAuth(base.Auth, override.Auth)
	if override.RegionName != "" {
		result.RegionName = override.RegionName
	}
	if override.Interface != "" {
		result.Interface = override.Interface
	}
	if override.IdentityAPIVersion != "" {
		result.IdentityAPIVersion = override.IdentityAPIVersion
	}
	if override.AuthType != "" {
		result.AuthType = override.AuthType
	}
	if override.ProjectDomainName != "" {
		result.ProjectDomainName = override.ProjectDomainName
	}
	return result
}

// mergeCloudAuth merges two cloudAuth structs. Any non-empty field in override takes priority.
func mergeCloudAuth(base, override cloudAuth) cloudAuth {
	result := base
	if override.AuthURL != "" {
		result.AuthURL = override.AuthURL
	}
	if override.Username != "" {
		result.Username = override.Username
	}
	if override.Password != "" {
		result.Password = override.Password
	}
	if override.ApplicationCredentialID != "" {
		result.ApplicationCredentialID = override.ApplicationCredentialID
	}
	if override.ApplicationCredentialSecret != "" {
		result.ApplicationCredentialSecret = override.ApplicationCredentialSecret
	}
	if override.ProjectName != "" {
		result.ProjectName = override.ProjectName
	}
	if override.ProjectID != "" {
		result.ProjectID = override.ProjectID
	}
	if override.UserDomainName != "" {
		result.UserDomainName = override.UserDomainName
	}
	if override.ProjectDomainName != "" {
		result.ProjectDomainName = override.ProjectDomainName
	}
	if override.ProjectDomainID != "" {
		result.ProjectDomainID = override.ProjectDomainID
	}
	return result
}

type cloudsYAML struct {
	Clouds map[string]cloudEntry `yaml:"clouds"`
}

type cloudEntry struct {
	Auth               cloudAuth `yaml:"auth"`
	RegionName         string    `yaml:"region_name"`
	Interface          string    `yaml:"interface"`
	IdentityAPIVersion string    `yaml:"identity_api_version"`
	AuthType           string    `yaml:"auth_type"`
	ProjectDomainName  string    `yaml:"project_domain_name"` // can also appear at entry level
}

type cloudAuth struct {
	AuthURL                     string `yaml:"auth_url"`
	Username                    string `yaml:"username"`
	Password                    string `yaml:"password"`
	ApplicationCredentialID     string `yaml:"application_credential_id"`
	ApplicationCredentialSecret string `yaml:"application_credential_secret"`
	ProjectName                 string `yaml:"project_name"`
	ProjectID                   string `yaml:"project_id"`
	UserDomainName              string `yaml:"user_domain_name"`
	ProjectDomainName           string `yaml:"project_domain_name"`
	ProjectDomainID             string `yaml:"project_domain_id"`
}
