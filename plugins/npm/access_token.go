package npm

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const defaultRegistry = "https://registry.npmjs.org/"

func AccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessToken,
		DocsURL:       sdk.URL("https://docs.npmjs.com/about-access-tokens"),
		ManagementURL: sdk.URL("https://www.npmjs.com/settings/~/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Access token used to authenticate to an npm-compatible registry.",
				Secret:              true,
			},
			{
				Name:                fieldname.Organization,
				MarkdownDescription: "The package scope this registry should be used for, without the leading @.",
				Optional:            true,
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "The npm-compatible registry host or URL that accepts this access token.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.TempFile(
			npmConfigFile,
			provision.Filename(".npmrc"),
			provision.AddArgs("--userconfig", "{{ .Path }}"),
		),
		Importer: importer.TryAll(
			importer.TryAllEnvVars(fieldname.Token, "NPM_TOKEN", "NODE_AUTH_TOKEN"),
			tryNPMRCFile("~/.npmrc"),
			tryPNPMAuthFile(),
			// pnpm 10 and earlier may have stored npm-compatible settings here.
			tryNPMRCFile("~/.config/pnpm/rc"),
		),
	}
}

func pnpmProvisioner() sdk.Provisioner {
	return provision.TempFile(
		npmConfigFile,
		provision.Filename(".npmrc"),
		provision.SetPathAsEnvVar("NPM_CONFIG_USERCONFIG"),
	)
}

func npmConfigFile(in sdk.ProvisionInput) ([]byte, error) {
	registry, err := normalizeRegistry(in.ItemFields[fieldname.Host])
	if err != nil {
		return nil, err
	}

	scope := strings.TrimPrefix(strings.TrimSpace(in.ItemFields[fieldname.Organization]), "@")
	var contents strings.Builder
	if scope != "" {
		fmt.Fprintf(&contents, "@%s:registry=%s\n", scope, registry.String())
	} else if strings.TrimSpace(in.ItemFields[fieldname.Host]) != "" {
		fmt.Fprintf(&contents, "registry=%s\n", registry.String())
	}

	registryPath := registry.EscapedPath()
	if !strings.HasSuffix(registryPath, "/") {
		registryPath += "/"
	}
	fmt.Fprintf(&contents, "//%s%s:_authToken=%s\n", registry.Host, registryPath, in.ItemFields[fieldname.Token])

	return []byte(contents.String()), nil
}

func normalizeRegistry(value string) (*url.URL, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		value = defaultRegistry
	} else if !strings.Contains(value, "://") {
		value = "https://" + value
	}

	registry, err := url.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("parsing registry URL: %w", err)
	}
	if (registry.Scheme != "https" && registry.Scheme != "http") || registry.Host == "" {
		return nil, fmt.Errorf("registry URL must use http or https and include a host")
	}
	if registry.User != nil || registry.RawQuery != "" || registry.Fragment != "" {
		return nil, fmt.Errorf("registry URL must not contain credentials, a query, or a fragment")
	}
	if registry.Path == "" {
		registry.Path = "/"
	} else if !strings.HasSuffix(registry.Path, "/") {
		registry.Path += "/"
	}

	return registry, nil
}

func tryPNPMAuthFile() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		var path string
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			path = filepath.Join(xdgConfigHome, "pnpm", "auth.ini")
		} else {
			switch in.OS {
			case "darwin":
				path = "~/Library/Preferences/pnpm/auth.ini"
			case "linux":
				path = "~/.config/pnpm/auth.ini"
			default:
				return
			}
		}

		tryNPMRCFile(path)(ctx, in, out)
	}
}

func tryNPMRCFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		lines := make(map[string]string)
		scanner := bufio.NewScanner(strings.NewReader(string(contents)))
		for scanner.Scan() {
			key, value, found := strings.Cut(scanner.Text(), "=")
			if !found {
				continue
			}
			lines[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
		}
		if err := scanner.Err(); err != nil {
			out.AddError(err)
			return
		}

		scopesByRegistry := make(map[string][]string)
		registryURLs := make(map[string]string)
		for key, value := range lines {
			lowerKey := strings.ToLower(key)
			if lowerKey != "registry" && (!strings.HasPrefix(key, "@") || !strings.HasSuffix(lowerKey, ":registry")) {
				continue
			}
			registry, err := normalizeRegistry(value)
			if err != nil {
				continue
			}
			registryKey := registryAuthKey(registry)
			registryURLs[registryKey] = registry.String()
			if lowerKey == "registry" {
				continue
			}
			scope := key[1 : len(key)-len(":registry")]
			if scope != "" {
				scopesByRegistry[registryKey] = append(scopesByRegistry[registryKey], scope)
			}
		}

		for key, token := range lines {
			registryKey, explicitScope, ok := parseAuthKey(key)
			if !ok || token == "" || strings.HasPrefix(token, "${") {
				continue
			}

			scopes := []string{explicitScope}
			if explicitScope == "" {
				if configuredScopes := scopesByRegistry[registryKey]; len(configuredScopes) > 0 {
					scopes = configuredScopes
				}
			}

			for _, scope := range scopes {
				host := registryKey
				if configuredURL := registryURLs[registryKey]; configuredURL != "" {
					host = configuredURL
				}
				fields := map[sdk.FieldName]string{
					fieldname.Token: token,
					fieldname.Host:  host,
				}
				if scope != "" {
					fields[fieldname.Organization] = scope
				}
				out.AddCandidate(sdk.ImportCandidate{
					Fields:   fields,
					NameHint: importer.SanitizeNameHint(registryKey),
				})
			}
		}
	})
}

func parseAuthKey(key string) (registry string, scope string, ok bool) {
	const suffix = ":_authtoken"
	lowerKey := strings.ToLower(key)
	if !strings.HasPrefix(key, "//") || !strings.HasSuffix(lowerKey, suffix) {
		return "", "", false
	}

	registry = key[2 : len(key)-len(suffix)]
	if scopeSeparator := strings.LastIndex(registry, ":@"); scopeSeparator >= 0 {
		scope = strings.TrimPrefix(registry[scopeSeparator+1:], "@")
		registry = registry[:scopeSeparator]
	}
	registry = strings.TrimSuffix(registry, ":")
	registry = strings.TrimSuffix(registry, "/")
	if registry == "" {
		return "", "", false
	}
	return registry, scope, true
}

func registryAuthKey(registry *url.URL) string {
	path := strings.TrimSuffix(registry.EscapedPath(), "/")
	return registry.Host + path
}
