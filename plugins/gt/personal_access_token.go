package gt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PersonalAccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://docs.gitea.com/development/api-usage/#generating-and-listing-api-tokens"),
		ManagementURL: sdk.URL("https://gitea.com/user/settings/applications"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Gitea.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.HostAddress,
				MarkdownDescription: "The Gitea host to connect to, should start with https://",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Symbols:   true,
					},
					Prefix: "https://",
				},
			},
			{
				Name:                fieldname.User,
				MarkdownDescription: "Your Gitea username",
				Secret:              false,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"GITEA_TOKEN":    fieldname.Token,
				"GITEA_CLI_HOST": fieldname.HostAddress,
			}),
			TryGiteaHostsFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"GITEA_TOKEN":        fieldname.Token,
	"GITEA_INSTANCE_URL": fieldname.HostAddress,
}

func configDir() string {
	if gitConfigCliDir := os.Getenv("GITEA_CLI_CONFIG_DIR"); gitConfigCliDir != "" {
		return gitConfigCliDir
	}

	// gt's config resolver honors XDG_CONFIG_HOME on macOS; os.UserConfigDir does not.
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "gitea-cli")
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "~/.config/gitea-cli"
	}
	return filepath.Join(configDir, "gitea-cli")
}

func HostsPath() string {
	return filepath.Join(configDir(), "hosts.yml")
}

// Identity builds a stable, opaque key from an instance URL and account.
func identity(instance, user string) string {
	sum := sha256.Sum256([]byte(instance + "\x00" + user))
	return hex.EncodeToString(sum[:])
}

func credentialFilename(instance string, user string) string {
	return filepath.Join(configDir(), "credentials", identity(instance, user)+".json")
}

func loadCredentialFromFile(hostname string, username string, rootDir string) (Credential, error) {
	var credential Credential
	// In test environments, absolute paths are prefixed with RootDir
	credPath := credentialFilename(hostname, username)
	if rootDir != "" && filepath.IsAbs(credPath) {
		credPath = filepath.Join(rootDir, credPath)
	}
	credBytes, err := os.ReadFile(credPath)
	if err != nil {
		return Credential{}, err
	}

	if err := json.Unmarshal(credBytes, &credential); err != nil {
		return Credential{}, err
	}

	return credential, nil
}

func TryGiteaHostsFile() sdk.Importer {
	return importer.TryFile(HostsPath(), func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var hosts Hosts
		if err := contents.ToYAML(&hosts); err != nil {
			out.AddError(err)
			return
		}

		for hostname, host := range hosts.Hosts {
			for username, account := range host.Users {
				if account.CredentialStore == "file" && account.AuthMethod == "token" {
					credential, err := loadCredentialFromFile(hostname, username, in.RootDir)
					if err != nil {
						out.AddError(err)
						continue
					}
					out.AddCandidate(sdk.ImportCandidate{
						Fields: map[sdk.FieldName]string{
							fieldname.Token:       credential.AccessToken,
							fieldname.HostAddress: hostname,
							fieldname.User:        username,
						},
						NameHint: importer.SanitizeNameHint(username + "@" + hostname),
					})
				}
			}
		}
	})
}

type Credential struct {
	AccessToken string `json:"access_token"`
}

type Account struct {
	AuthMethod      string `yaml:"auth_method" json:"auth_method"`
	CredentialStore string `yaml:"credential_store,omitempty" json:"credential_store,omitempty"`
}

type Host struct {
	User  string             `yaml:"user,omitempty" json:"user"`
	Users map[string]Account `yaml:"users" json:"users"`
}

type Hosts struct {
	Version string          `yaml:"version" json:"version"`
	Hosts   map[string]Host `yaml:"hosts" json:"hosts"`
}
